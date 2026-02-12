import { create } from "zustand";
import { userApi } from "../services/api";
import { useAuthStore } from "./authStore";

export interface Meal {
  id: string;
  user_id: string;
  timestamp: string; // ISO date string
  name: string;
  description?: string;
  image_path?: string;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  confidence_level?: number;
}

interface DailyStats {
  totalCalories: number;
  totalProtein: number;
  totalCarbs: number;
  totalFat: number;
}

interface MealsState {
  meals: Meal[];
  selectedDate: string; // YYYY-MM-DD
  isLoading: boolean;
  error: string | null;

  // Actions
  setDate: (date: string) => void;
  fetchMeals: () => Promise<void>;
  addMeal: (meal: Omit<Meal, "id" | "timestamp" | "user_id">) => Promise<void>;

  // Selectors (computed values)
  getDailyStats: () => DailyStats;
}

export const useMealsStore = create<MealsState>((set, get) => ({
  meals: [],
  selectedDate: new Date().toISOString().split("T")[0],
  isLoading: false,
  error: null,

  setDate: (date: string) => {
    set({ selectedDate: date });
    get().fetchMeals();
  },

  fetchMeals: async () => {
    const { selectedDate } = get();
    const user = useAuthStore.getState().user;
    if (!user) return;

    set({ isLoading: true, error: null });
    try {
      const response = await userApi.getMeals(user.id, selectedDate);
      set({ meals: response.data || [], isLoading: false });
    } catch (error: any) {
      console.error("Fetch meals error:", error);
      set({ error: "Failed to fetch meals", isLoading: false });
    }
  },

  addMeal: async (mealData) => {
    const user = useAuthStore.getState().user;
    if (!user) return;

    set({ isLoading: true, error: null });
    try {
      // Optimistic update could go here, but let's just refresh for now to be safe with DB IDs
      const newMeal = {
        ...mealData,
        user_id: user.id,
        // timestamp handled by backend or we send it? Backend creates it usually.
        // API requirement might need checking.
        // Let's assume we send current time if creating for "now", but for a specific date?
        // The prompt says "Log Meal", usually implies "now".
      };

      await userApi.logMeal(user.id, newMeal);

      // Refresh list
      await get().fetchMeals();
      set({ isLoading: false });
    } catch (error: any) {
      console.error("Add meal error:", error);
      set({ error: "Failed to log meal", isLoading: false });
      throw error; // Re-throw so UI can handle navigation
    }
  },

  getDailyStats: () => {
    const { meals } = get();
    return (meals || []).reduce(
      (acc, meal) => ({
        totalCalories: acc.totalCalories + meal.calories,
        totalProtein: acc.totalProtein + meal.protein_g,
        totalCarbs: acc.totalCarbs + meal.carbs_g,
        totalFat: acc.totalFat + meal.fat_g,
      }),
      { totalCalories: 0, totalProtein: 0, totalCarbs: 0, totalFat: 0 },
    );
  },
}));
