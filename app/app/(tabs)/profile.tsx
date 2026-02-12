import { useState, useEffect } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from "react-native";
import { useAuthStore } from "../../src/store/authStore";
import { useRouter, useFocusEffect } from "expo-router";
import { userApi } from "../../src/services/api";
import { CreateMacrosRequest, Macros } from "../../src/types/user";
import React from "react";

export default function ProfileScreen() {
  const { user, logout } = useAuthStore();
  const router = useRouter();

  console.log("ProfileScreen: Rendering Full Component"); // Debug log

  // Password State
  const [newPassword, setNewPassword] = useState("");
  const [loadingPassword, setLoadingPassword] = useState(false);

  // Macros State
  const [macrosList, setMacrosList] = useState<Macros[]>([]);
  const [activeMacro, setActiveMacro] = useState<Macros | null>(null);
  const [mode, setMode] = useState<"view" | "create" | "history">("create");

  // Form State
  const [targetCalories, setTargetCalories] = useState("");
  const [targetProtein, setTargetProtein] = useState("");
  const [targetCarbs, setTargetCarbs] = useState("");
  const [targetFat, setTargetFat] = useState("");
  const [loadingMacros, setLoadingMacros] = useState(false);

  // Fetch macros when screen comes into focus or user changes
  useFocusEffect(
    React.useCallback(() => {
      if (user) {
        fetchMacros();
      }
    }, [user]),
  );

  const fetchMacros = async () => {
    if (!user) return;
    try {
      const response = await userApi.getMacros(user.id);
      console.log("ProfileScreen: Macros fetched", response.data);
      if (Array.isArray(response.data) && response.data.length > 0) {
        // Sort by timestamp desc to get history right
        const sorted = response.data.sort(
          (a, b) =>
            new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
        );
        setMacrosList(sorted);

        const active = sorted.find((m) => m.active) || sorted[0];
        setActiveMacro(active);

        // Only switch to view if we are not already in history or create
        // Actually, initial load should determine mode.
        // If we have active macro, default to view.
        setMode("view");
      } else {
        setMacrosList([]);
        setActiveMacro(null);
        setMode("create");
      }
    } catch (error) {
      console.error("Failed to fetch macros:", error);
    }
  };

  const fetchHistory = async () => {
    if (!user) return;
    try {
      const response = await userApi.getHistoricalMacros(user.id);
      if (Array.isArray(response.data)) {
        const sorted = response.data.sort(
          (a, b) =>
            new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
        );
        setMacrosList(sorted);
      }
    } catch (error) {
      console.error("Failed to fetch macro history:", error);
    }
  };

  const handleLogout = async () => {
    router.replace("/(auth)/logout");
  };

  const handleChangePassword = async () => {
    if (!user) return;
    if (!newPassword.trim()) {
      Alert.alert("Error", "Please enter a new password");
      return;
    }

    setLoadingPassword(true);
    try {
      await userApi.changePassword(user.id, newPassword);
      Alert.alert("Success", "Password changed successfully");
      setNewPassword("");
    } catch (error: any) {
      console.error("Change password error:", error);
      let msg = "Failed to change password";
      if (error.response?.data?.error) msg = error.response.data.error;
      Alert.alert("Error", msg);
    } finally {
      setLoadingPassword(false);
    }
  };

  const handleSaveMacros = async () => {
    if (!user) return;

    const cal = parseInt(targetCalories);
    const prot = parseFloat(targetProtein);
    const carbs = parseFloat(targetCarbs);
    const fat = parseFloat(targetFat);

    if (isNaN(cal) || isNaN(prot) || isNaN(carbs) || isNaN(fat)) {
      Alert.alert("Error", "Please enter valid numbers for all macro targets");
      return;
    }

    setLoadingMacros(true);
    try {
      const payload: CreateMacrosRequest = {
        calories: cal,
        protein_g: prot,
        carbs_g: carbs,
        fat_g: fat,
      };

      // Always create new entry
      await userApi.addMacros(user.id, payload);
      Alert.alert("Success", "New macros created successfully");

      // Clear form
      setTargetCalories("");
      setTargetProtein("");
      setTargetCarbs("");
      setTargetFat("");

      // Refresh
      await fetchMacros();
    } catch (error: any) {
      console.error("Save macros error:", error);
      let msg = "Failed to save macros";
      if (error.response?.data?.error) msg = error.response.data.error;
      Alert.alert("Error", msg);
    } finally {
      setLoadingMacros(false);
    }
  };

  // Render Helpers
  const renderMacroView = () => {
    if (!activeMacro)
      return (
        <View>
          <Text style={styles.infoText}>No active macros found.</Text>
          <TouchableOpacity
            style={styles.primaryButton}
            onPress={() => setMode("create")}
          >
            <Text style={styles.primaryButtonText}>Create New Targets</Text>
          </TouchableOpacity>
        </View>
      );

    return (
      <View>
        <View style={styles.macroCard}>
          <Text style={styles.macroTitle}>current Active Targets</Text>
          <Text style={styles.macroDate}>
            Set on: {new Date(activeMacro.timestamp).toLocaleDateString()}
          </Text>

          <View style={styles.macroGrid}>
            <View style={styles.macroItem}>
              <Text style={styles.macroValueBig}>{activeMacro.calories}</Text>
              <Text style={styles.macroLabel}>kcal</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValueBig}>{activeMacro.protein_g}</Text>
              <Text style={styles.macroLabel}>Prot (g)</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValueBig}>{activeMacro.carbs_g}</Text>
              <Text style={styles.macroLabel}>Carbs (g)</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValueBig}>{activeMacro.fat_g}</Text>
              <Text style={styles.macroLabel}>Fat (g)</Text>
            </View>
          </View>
        </View>

        <View style={styles.buttonRow}>
          <TouchableOpacity
            style={styles.secondaryButton}
            onPress={() => {
              setMode("history");
              fetchHistory();
            }}
          >
            <Text style={styles.secondaryButtonText}>History</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={styles.primaryButton}
            onPress={() => setMode("create")}
          >
            <Text style={styles.primaryButtonText}>Update</Text>
          </TouchableOpacity>
        </View>
      </View>
    );
  };

  const renderMacroForm = () => (
    <View>
      <Text style={styles.subHeader}>Set New Macro Targets</Text>
      <Text style={styles.label}>Calories (kcal)</Text>
      <TextInput
        style={styles.input}
        placeholder="e.g. 2500"
        placeholderTextColor="#666"
        keyboardType="numeric"
        value={targetCalories}
        onChangeText={setTargetCalories}
      />

      <View style={styles.row}>
        <View style={styles.col}>
          <Text style={styles.label}>Protein (g)</Text>
          <TextInput
            style={styles.input}
            placeholder="180"
            placeholderTextColor="#666"
            keyboardType="numeric"
            value={targetProtein}
            onChangeText={setTargetProtein}
          />
        </View>
        <View style={[styles.col, styles.marginH]}>
          <Text style={styles.label}>Carbs (g)</Text>
          <TextInput
            style={styles.input}
            placeholder="250"
            placeholderTextColor="#666"
            keyboardType="numeric"
            value={targetCarbs}
            onChangeText={setTargetCarbs}
          />
        </View>
        <View style={styles.col}>
          <Text style={styles.label}>Fat (g)</Text>
          <TextInput
            style={styles.input}
            placeholder="80"
            placeholderTextColor="#666"
            keyboardType="numeric"
            value={targetFat}
            onChangeText={setTargetFat}
          />
        </View>
      </View>

      <View style={styles.buttonRow}>
        {macrosList.length > 0 && (
          <TouchableOpacity
            style={styles.secondaryButton}
            onPress={() => setMode("view")}
          >
            <Text style={styles.secondaryButtonText}>Cancel</Text>
          </TouchableOpacity>
        )}
        <TouchableOpacity
          style={[
            styles.saveButton,
            loadingMacros && styles.disabledButton,
            { flex: 1, marginLeft: 8 },
          ]}
          onPress={handleSaveMacros}
          disabled={loadingMacros}
        >
          {loadingMacros ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.buttonText}>Save New Targets</Text>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );

  const renderHistory = () => (
    <View>
      <Text style={styles.subHeader}>Macro History</Text>
      {macrosList.map((m) => (
        <View
          key={m.id}
          style={[styles.historyItem, m.active && styles.activeHistoryItem]}
        >
          <View style={styles.historyRow}>
            <Text style={styles.historyDate}>
              {new Date(m.timestamp).toLocaleDateString()}{" "}
              {new Date(m.timestamp).toLocaleTimeString()}
            </Text>
            {m.active && <Text style={styles.activeBadge}>ACTIVE</Text>}
          </View>
          <Text style={styles.historyText}>
            {m.calories} kcal | P: {m.protein_g} | C: {m.carbs_g} | F: {m.fat_g}
          </Text>
        </View>
      ))}
      <TouchableOpacity
        style={[styles.secondaryButton, { marginTop: 16 }]}
        onPress={() => {
          setMode("view");
          fetchMacros();
        }}
      >
        <Text style={styles.secondaryButtonText}>Back to Active</Text>
      </TouchableOpacity>
    </View>
  );

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      keyboardVerticalOffset={Platform.OS === "ios" ? 100 : 0}
      style={styles.container}
    >
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* SECTION A: USER INFO */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Profile Info</Text>
          <View style={styles.infoRow}>
            <Text style={styles.label}>Name:</Text>
            <Text style={styles.value}>{user?.name}</Text>
          </View>
          <View style={styles.infoRow}>
            <Text style={styles.label}>Email:</Text>
            <Text style={styles.value}>{user?.email}</Text>
          </View>

          <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
            <Text style={styles.logoutButtonText}>Logout</Text>
          </TouchableOpacity>
        </View>

        {/* SECTION B: CHANGE PASSWORD */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Change Password</Text>
          <TextInput
            style={styles.input}
            placeholder="New Password"
            placeholderTextColor="#666"
            secureTextEntry={true}
            value={newPassword}
            onChangeText={setNewPassword}
          />
          <TouchableOpacity
            style={[
              styles.actionButton,
              loadingPassword && styles.disabledButton,
            ]}
            onPress={handleChangePassword}
            disabled={loadingPassword}
          >
            {loadingPassword ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.buttonText}>Update Password</Text>
            )}
          </TouchableOpacity>
        </View>

        {/* SECTION C: MACROS UI */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Macro Targets</Text>
          {mode === "view" && renderMacroView()}
          {mode === "create" && renderMacroForm()}
          {mode === "history" && renderHistory()}
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#121212",
  },
  scrollContent: {
    padding: 20,
    paddingBottom: 40,
  },
  section: {
    backgroundColor: "#1E1E1E",
    borderRadius: 12,
    padding: 20,
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: "bold",
    color: "#FFFFFF",
    marginBottom: 16,
    borderBottomWidth: 1,
    borderBottomColor: "#333",
    paddingBottom: 8,
  },
  infoRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 12,
  },
  label: {
    color: "#AAAAAA",
    fontSize: 14,
    marginBottom: 6,
  },
  value: {
    color: "#FFFFFF",
    fontSize: 16,
    fontWeight: "600",
  },
  input: {
    backgroundColor: "#2C2C2C",
    color: "#FFFFFF",
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: "#333",
  },
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
  },
  col: {
    flex: 1,
  },
  marginH: {
    marginHorizontal: 10,
  },
  logoutButton: {
    backgroundColor: "#FF3B30",
    padding: 14,
    borderRadius: 8,
    alignItems: "center",
    marginTop: 10,
  },
  logoutButtonText: {
    color: "#FFFFFF",
    fontWeight: "bold",
    fontSize: 16,
  },
  adminButton: {
    backgroundColor: "#5856D6",
    padding: 14,
    borderRadius: 8,
    alignItems: "center",
    marginTop: 10,
    marginBottom: 8,
  },
  adminButtonText: {
    color: "#FFFFFF",
    fontWeight: "bold",
    fontSize: 16,
  },
  actionButton: {
    backgroundColor: "#007AFF",
    padding: 14,
    borderRadius: 8,
    alignItems: "center",
  },
  saveButton: {
    backgroundColor: "#34C759",
    padding: 14,
    borderRadius: 8,
    alignItems: "center",
  },
  buttonText: {
    color: "#FFFFFF",
    fontWeight: "bold",
    fontSize: 16,
  },
  disabledButton: {
    opacity: 0.6,
  },
  // New Styles
  buttonRow: {
    flexDirection: "row",
    marginTop: 16,
  },
  primaryButton: {
    backgroundColor: "#34C759",
    padding: 12,
    borderRadius: 8,
    flex: 1,
    marginLeft: 8,
    alignItems: "center",
  },
  primaryButtonText: {
    color: "#FFF",
    fontWeight: "bold",
  },
  secondaryButton: {
    backgroundColor: "#3A3A3C",
    padding: 12,
    borderRadius: 8,
    flex: 1,
    marginRight: 8,
    alignItems: "center",
  },
  secondaryButtonText: {
    color: "#FFF",
    fontWeight: "600",
  },
  subHeader: {
    fontSize: 16, // Smaller than section title
    color: "#DDD",
    marginBottom: 16,
    fontWeight: "600",
  },
  macroCard: {
    backgroundColor: "#252525",
    padding: 16,
    borderRadius: 8,
    marginBottom: 8,
    alignItems: "center",
  },
  macroTitle: {
    color: "#34C759",
    fontWeight: "bold",
    marginBottom: 4,
    fontSize: 14,
    textTransform: "uppercase",
    letterSpacing: 1,
  },
  macroDate: {
    color: "#666",
    fontSize: 12,
    marginBottom: 16,
  },
  macroGrid: {
    flexDirection: "row",
    justifyContent: "space-between",
    width: "100%",
  },
  macroItem: {
    alignItems: "center",
    flex: 1,
  },
  macroValueBig: {
    color: "#FFF",
    fontSize: 20,
    fontWeight: "bold",
  },
  macroLabel: {
    color: "#AAA",
    fontSize: 12,
    marginTop: 2,
  },
  infoText: {
    color: "#888",
    textAlign: "center",
    marginBottom: 16,
    fontStyle: "italic",
  },
  historyItem: {
    backgroundColor: "#2C2C2C",
    padding: 12,
    borderRadius: 8,
    marginBottom: 8,
    borderLeftWidth: 4,
    borderLeftColor: "#444",
  },
  activeHistoryItem: {
    borderLeftColor: "#34C759",
    backgroundColor: "#253525",
  },
  historyRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 4,
  },
  historyDate: {
    color: "#FFF",
    fontWeight: "bold",
    fontSize: 14,
  },
  activeBadge: {
    color: "#34C759",
    fontSize: 10,
    fontWeight: "bold",
    borderWidth: 1,
    borderColor: "#34C759",
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
  historyText: {
    color: "#CCC",
    fontSize: 13,
  },
});
