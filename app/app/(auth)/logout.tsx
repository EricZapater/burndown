import { useEffect } from "react";
import { View, Text, StyleSheet, ActivityIndicator } from "react-native";
import { useRouter } from "expo-router";
import { useAuthStore } from "../../src/store/authStore";

export default function LogoutScreen() {
  const { logout, isAuthenticated } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    // Perform logout
    logout();
  }, []);

  useEffect(() => {
    if (!isAuthenticated) {
      // Optional delay for user to see the message
      const timer = setTimeout(() => {
        router.replace("/");
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [isAuthenticated]);

  return (
    <View style={styles.container}>
      <Text style={styles.text}>Logging out...</Text>
      <Text style={styles.subText}>See you next time!</Text>
      <ActivityIndicator size="large" color="#34C759" style={styles.loader} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#121212",
    justifyContent: "center",
    alignItems: "center",
  },
  text: {
    color: "#fff",
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 8,
  },
  subText: {
    color: "#888",
    fontSize: 16,
  },
  loader: {
    marginTop: 20,
  },
});
