import { useEffect, useState } from "react";
import {
  View,
  Text,
  FlatList,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useAuthStore } from "../../src/store/authStore";
import { userApi } from "../../src/services/api";

export default function UsersListScreen() {
  const router = useRouter();
  const { user } = useAuthStore();
  const [users, setUsers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Double check just in case, though _layout should handle visibility
    if (!user?.is_admin) {
      // If they somehow got here, just show nothing or redirect
      return;
    }
    fetchUsers();
  }, [user]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const res = await userApi.getUsers();
      setUsers(res.data);
    } catch (error) {
      console.error("Failed to fetch users", error);
      Alert.alert("Error", "Failed to fetch users.");
    } finally {
      setLoading(false);
    }
  };

  const renderItem = ({ item }: { item: any }) => (
    <View style={styles.card}>
      <View style={styles.info}>
        <Text style={styles.name}>{item.name}</Text>
        <Text style={styles.email}>{item.email}</Text>
        <View style={styles.tags}>
          {item.is_admin ? <Text style={styles.adminTag}>ADMIN</Text> : null}
          {!item.active ? (
            <Text style={styles.inactiveTag}>INACTIVE</Text>
          ) : null}
        </View>
      </View>
    </View>
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        {/* No back button for main tab */}
        <Text style={styles.title}>User Management</Text>
        <TouchableOpacity
          onPress={() => router.push("/admin/create-user")}
          style={styles.addButton}
        >
          <Ionicons name="add" size={24} color="#FFF" />
        </TouchableOpacity>
      </View>

      {loading ? (
        <ActivityIndicator
          size="large"
          color="#34C759"
          style={{ marginTop: 20 }}
        />
      ) : (
        <FlatList
          data={users}
          keyExtractor={(item) => item.id}
          renderItem={renderItem}
          contentContainerStyle={styles.list}
          refreshing={loading}
          onRefresh={fetchUsers}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#121212" },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    padding: 16,
    paddingTop: 60, // Extra padding for status bar since headerShown: false
    backgroundColor: "#1E1E1E",
  },
  addButton: { padding: 8 },
  title: { color: "#FFF", fontSize: 20, fontWeight: "bold" },
  list: { padding: 16 },
  card: {
    backgroundColor: "#1E1E1E",
    padding: 16,
    borderRadius: 8,
    marginBottom: 12,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  info: { flex: 1 },
  name: { color: "#FFF", fontSize: 16, fontWeight: "bold" },
  email: { color: "#AAA", fontSize: 14, marginTop: 4 },
  tags: { flexDirection: "row", marginTop: 8 },
  adminTag: {
    backgroundColor: "#5856D6",
    color: "#FFF",
    fontSize: 10,
    fontWeight: "bold",
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
    marginRight: 8,
  },
  inactiveTag: {
    backgroundColor: "#FF3B30",
    color: "#FFF",
    fontSize: 10,
    fontWeight: "bold",
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 4,
  },
});
