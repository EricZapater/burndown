import axios from "axios";
import { useAuthStore } from "../store/authStore";
import { CreateMacrosRequest, Macros } from "../types/user";

// ⚠️ TODO: CHANGE THIS IP (192.168.1.X) TO YOUR COMPUTER'S LOCAL IP
// Example: If your computer's IP is 192.168.1.100, use http://192.168.1.100:8080
// Find your IP: Windows CMD -> ipconfig (look for IPv4 Address under your network adapter)
//const API_URL = "http://192.168.68.105:8125/api/v1";
const API_URL = process.env.EXPO_PUBLIC_API_URL;

const api = axios.create({
  baseURL: API_URL,
  timeout: 60000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to inject JWT token
api.interceptors.request.use(
  (config) => {
    console.log(
      `API Request: ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`,
    );
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    console.error("API Request Error:", error);
    return Promise.reject(error);
  },
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error(
      `API Response Error: ${error.message} for ${error.config?.url}`,
    );
    if (error.response?.status === 401) {
      // Token expired or invalid - logout
      useAuthStore.getState().logout();
    }
    return Promise.reject(error);
  },
);

export const userApi = {
  changePassword: async (userId: string, password: string) => {
    // Backend expects a raw string or single-field JSON, trying raw string based on Handler analysis
    return api.put(`/users/${userId}/password`, JSON.stringify(password));
  },

  getMacros: async (userId: string) => {
    return api.get<Macros[]>(`/users/${userId}/macros`);
  },

  getHistoricalMacros: async (userId: string) => {
    return api.get<Macros[]>(`/users/${userId}/macros/historical`);
  },

  addMacros: async (userId: string, macros: CreateMacrosRequest) => {
    return api.post(`/users/${userId}/macros`, macros);
  },

  updateMacros: async (
    userId: string,
    macroId: string,
    macros: CreateMacrosRequest,
  ) => {
    return api.put(`/users/${userId}/macros/${macroId}`, macros);
  },

  analyzeImage: async (imageUri: string, prompt?: string) => {
    const formData = new FormData();
    formData.append("image", {
      uri: imageUri,
      name: "meal.jpg",
      type: "image/jpeg",
    } as any);

    if (prompt) {
      formData.append("prompt", prompt);
    }

    // Axios requires specific headers for FormData in React Native
    return api.post("/advisor/analyze", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      transformRequest: (data, headers) => {
        // Axios hack for React Native FormData
        return formData;
      },
    });
  },

  getMeals: async (userId: string, date: string) => {
    // Matches the new route: users.GET("/:id/meals/date/:date", handler.FindAllByUserIdAndDate)
    return api.get(`/users/${userId}/meals/date/${date}`);
  },

  logMeal: async (userId: string, mealData: any) => {
    // Ensure the structure matches CreateMealRequest in backend
    // UserID is in the URL in some designs but body in others.
    // Backend 'Create' handler binds JSON. 'CreateMealRequest' has 'user_id'.
    // Handler is `meals.POST("", handler.Create)`.
    // So we post to `/meals` and body must contain `user_id`.

    // The store addMeal prepares mealData with user_id.
    // We need to pass that to the endpoint.

    return api.post(`/meals`, mealData);
  },

  getUsers: async () => {
    return api.get<any[]>("/users");
  },

  createUser: async (userData: any) => {
    return api.post("/users", userData); // Backend expects { name, email, password }
  },
};

export default api;
