import { useEffect, useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
  SafeAreaView,
  RefreshControl,
  StatusBar,
} from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useAuthStore } from "../../src/store/authStore";
import { useMealsStore, Meal } from "../../src/store/mealsStore";
import { userApi } from "../../src/services/api";

// Helper for date formatting
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
};

const getRelativeDateLabel = (dateString: string) => {
  const today = new Date().toISOString().split("T")[0];
  if (dateString === today) return "Today";

  // Check yesterday
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  const yesterdayStr = yesterday.toISOString().split("T")[0];
  if (dateString === yesterdayStr) return "Yesterday";

  return formatDate(dateString);
};

export default function DashboardScreen() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const { meals, selectedDate, isLoading, fetchMeals, setDate, getDailyStats } =
    useMealsStore();

  const [refreshing, setRefreshing] = useState(false);
  const [macroTargets, setMacroTargets] = useState<any>(null);

  useEffect(() => {
    fetchMeals();
    fetchUserMacros();
  }, [selectedDate]); // Refresh when date changes (meals) - macros maybe less often but this is fine

  const fetchUserMacros = async () => {
    if (!user) return;
    try {
      const res = await userApi.getMacros(user.id);
      if (res.data && res.data.length > 0) {
        const active = res.data.find((m: any) => m.active) || res.data[0];
        setMacroTargets(active);
      }
    } catch (e) {
      console.error("Failed to fetch macros for dashboard", e);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await Promise.all([fetchMeals(), fetchUserMacros()]);
    setRefreshing(false);
  };

  const handleDateChange = (days: number) => {
    const date = new Date(selectedDate);
    date.setDate(date.getDate() + days);
    setDate(date.toISOString().split("T")[0]);
  };

  const stats = getDailyStats();

  // Calculate Progress (capped at 100%)
  const getProgress = (current: number, target: number) => {
    if (!target || target === 0) return 0;
    return Math.min(current / target, 1);
  };

  const targets = {
    calories: macroTargets?.calories || user?.targetCalories || 2000,
    protein: macroTargets?.protein_g || user?.targetProtein || 150,
    carbs: macroTargets?.carbs_g || user?.targetCarbs || 200,
    fat: macroTargets?.fat_g || user?.targetFat || 70,
  };

  const caloriesProgress = getProgress(stats.totalCalories, targets.calories);
  const proteinProgress = getProgress(stats.totalProtein, targets.protein);
  const carbsProgress = getProgress(stats.totalCarbs, targets.carbs);
  const fatProgress = getProgress(stats.totalFat, targets.fat);

  // Color logic for calories (Green -> Red if over)
  // For now simple reliable colors
  const caloriesColor =
    stats.totalCalories > targets.calories ? "#FF3B30" : "#34C759";

  const renderProgressBar = (
    label: string,
    current: number,
    target: number,
    progress: number,
    color: string,
  ) => (
    <View style={styles.progressRow}>
      <View style={styles.progressLabelContainer}>
        <Text style={styles.progressLabel}>{label}</Text>
        <Text style={styles.progressValue}>
          {Math.round(current)} / {target}
        </Text>
      </View>
      <View style={styles.progressBarBg}>
        <View
          style={[
            styles.progressBarFill,
            { width: `${progress * 100}%`, backgroundColor: color },
          ]}
        />
      </View>
    </View>
  );

  const renderMealItem = ({ item }: { item: Meal }) => (
    <View style={styles.mealItem}>
      <View style={styles.mealInfo}>
        <Text style={styles.mealTime}>
          {new Date(
            item.timestamp || new Date().toISOString(),
          ).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
        </Text>
        <Text style={styles.mealName}>{item.name}</Text>
      </View>
      <View style={styles.mealStats}>
        <Text style={styles.mealCalories}>{item.calories} kcal</Text>
        <Text style={styles.mealMacros}>
          P:{Math.round(item.protein_g)} C:{Math.round(item.carbs_g)} F:
          {Math.round(item.fat_g)}
        </Text>
      </View>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="light-content" />

      {/* Header / Date Selector */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => handleDateChange(-1)}>
          <Ionicons name="chevron-back" size={24} color="#FFF" />
        </TouchableOpacity>

        <View style={styles.dateContainer}>
          <Text style={styles.dateLabel}>
            {getRelativeDateLabel(selectedDate)}
          </Text>
          <Text style={styles.fullDate}>{formatDate(selectedDate)}</Text>
        </View>

        <TouchableOpacity onPress={() => handleDateChange(1)}>
          <Ionicons name="chevron-forward" size={24} color="#FFF" />
        </TouchableOpacity>
      </View>

      <FlatList
        data={meals}
        renderItem={renderMealItem}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor="#FFF"
          />
        }
        ListHeaderComponent={
          <View style={styles.summaryCard}>
            <Text style={styles.summaryTitle}>Daily Summary</Text>

            {/* Main Calories Circle or Bar */}
            <View style={styles.caloriesHero}>
              <Text style={styles.heroValue}>
                {Math.round(targets.calories - stats.totalCalories)}
              </Text>
              <Text style={styles.heroLabel}>Calories Remaining</Text>
            </View>

            {/* Progress Bars */}
            {renderProgressBar(
              "Calories",
              stats.totalCalories,
              targets.calories,
              caloriesProgress,
              caloriesColor,
            )}
            {renderProgressBar(
              "Protein",
              stats.totalProtein,
              targets.protein,
              proteinProgress,
              "#5856D6",
            )}
            {renderProgressBar(
              "Carbs",
              stats.totalCarbs,
              targets.carbs,
              carbsProgress,
              "#FF9500",
            )}
            {renderProgressBar(
              "Fat",
              stats.totalFat,
              targets.fat,
              fatProgress,
              "#FFCC00",
            )}
          </View>
        }
        ListEmptyComponent={
          !isLoading ? (
            <View style={styles.emptyContainer}>
              <Text style={styles.emptyText}>No meals logged today.</Text>
              <Text style={styles.emptySubText}>
                Tap the + button to scan a meal!
              </Text>
            </View>
          ) : (
            <ActivityIndicator
              size="large"
              color="#FFF"
              style={{ marginTop: 50 }}
            />
          )
        }
      />

      {/* FAB */}
      <TouchableOpacity
        style={styles.fab}
        onPress={() => router.push("/(tabs)/scan")}
        activeOpacity={0.8}
      >
        <Ionicons name="camera" size={30} color="#000" />
      </TouchableOpacity>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#121212" },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: 20,
    paddingVertical: 15,
    borderBottomWidth: 1,
    borderBottomColor: "#2C2C2C",
  },
  dateContainer: { alignItems: "center" },
  dateLabel: { color: "#FFF", fontSize: 18, fontWeight: "bold" },
  fullDate: { color: "#AAA", fontSize: 12 },

  listContent: { paddingBottom: 100 },

  summaryCard: {
    backgroundColor: "#1E1E1E",
    margin: 16,
    borderRadius: 16,
    padding: 20,
    marginBottom: 24,
  },
  summaryTitle: {
    color: "#FFF",
    fontSize: 18,
    fontWeight: "bold",
    marginBottom: 20,
  },

  caloriesHero: { alignItems: "center", marginBottom: 25 },
  heroValue: { color: "#FFF", fontSize: 42, fontWeight: "900" },
  heroLabel: { color: "#AAA", fontSize: 14, textTransform: "uppercase" },

  progressRow: { marginBottom: 15 },
  progressLabelContainer: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 6,
  },
  progressLabel: { color: "#CCC", fontSize: 14 },
  progressValue: { color: "#FFF", fontSize: 14, fontWeight: "bold" },
  progressBarBg: {
    height: 8,
    backgroundColor: "#333",
    borderRadius: 4,
    overflow: "hidden",
  },
  progressBarFill: { height: "100%", borderRadius: 4 },

  mealItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    backgroundColor: "#1E1E1E",
    padding: 16,
    marginHorizontal: 16,
    marginBottom: 12,
    borderRadius: 12,
    alignItems: "center",
  },
  mealInfo: { flex: 1 },
  mealTime: { color: "#AAA", fontSize: 12, marginBottom: 4 },
  mealName: { color: "#FFF", fontSize: 16, fontWeight: "600" },
  mealStats: { alignItems: "flex-end" },
  mealCalories: { color: "#34C759", fontSize: 16, fontWeight: "bold" },
  mealMacros: { color: "#AAA", fontSize: 10, marginTop: 2 },

  emptyContainer: { alignItems: "center", marginTop: 50 },
  emptyText: { color: "#FFF", fontSize: 18, fontWeight: "bold" },
  emptySubText: { color: "#666", marginTop: 8 },

  fab: {
    position: "absolute",
    bottom: 25,
    right: 25,
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: "#34C759",
    alignItems: "center",
    justifyContent: "center",
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 4,
    elevation: 5,
  },
});
