import { Stack, useRouter, useSegments } from "expo-router";
import { useEffect } from "react";
import { useAuthStore } from "../src/store/authStore";

export default function RootLayout() {
  const segments = useSegments();
  const router = useRouter();
  const { isAuthenticated, isLoading, hydrate } = useAuthStore();

  // Hydrate auth state on mount
  useEffect(() => {
    hydrate();
  }, []);

  // Handle navigation based on auth state
  useEffect(() => {
    if (isLoading) return;

    const inAuthGroup = segments[0] === "(auth)";
    const inTabsGroup = segments[0] === "(tabs)";
    const inAdminGroup = segments[0] === "admin";
    const atRoot = segments.length === 0;

    const isLogout = segments[1] === "logout"; // Check if we are in logout route inside (auth)

    console.log("RootLayout Check:", {
      isAuthenticated,
      segments,
      inAuthGroup,
      inTabsGroup,
      atRoot,
    });

    if (!isAuthenticated) {
      // Not authenticated
      // Allow access to (auth) group and root (login)
      if (!inAuthGroup && !atRoot) {
        console.log("Redirecting to login (/)");
        router.replace("/");
      }
    } else {
      // Authenticated
      // Redirect to tabs if trying to access public screens (root or auth)
      // BUT allow (tabs), admin, AND logout
      if (!inTabsGroup && !inAdminGroup && !isLogout) {
        console.log("Redirecting to tabs (/(tabs))");
        router.replace("/(tabs)");
      }
    }
  }, [isAuthenticated, isLoading, segments]);

  return (
    <Stack
      screenOptions={{
        headerShown: false, // ✅ Strict boolean
      }}
    >
      <Stack.Screen name="index" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
    </Stack>
  );
}
