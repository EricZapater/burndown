import { useState, useRef } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Image,
  TextInput,
  ActivityIndicator,
  Alert,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  SafeAreaView,
} from "react-native";
import { CameraView, useCameraPermissions } from "expo-camera";
import { Ionicons } from "@expo/vector-icons";
import { useRouter } from "expo-router";
import { userApi } from "../../src/services/api";
import { useMealsStore } from "../../src/store/mealsStore";

export default function ScanScreen() {
  const [permission, requestPermission] = useCameraPermissions();
  const cameraRef = useRef<CameraView>(null);
  const router = useRouter();
  const addMeal = useMealsStore((state) => state.addMeal);
  const [loadingMessage, setLoadingMessage] = useState(
    "Analyzing your meal...",
  );

  // Funny Loading Messages
  const FUNNY_MESSAGES = [
    "La IA s'està posant el pitet... 🤤",
    "Comptant els pèsols un per un... 🔍",
    "Buscant la proteïna amagada... 🥩",
    "Negociant amb les calories perquè baixin... 📉",
    "Això va directe als bíceps, oi? 💪",
    "Calculant el nivell de penediment post-àpat... 😅",
    "La IA s'està menjant la foto amb els ulls... 👀",
    "Mesurant els macros amb regle i cartabó... 📐",
    "Tranquil, si t'ho menges ràpid la IA no ho veu... 🏃‍♂️",
    "Processant... si us plau, no et mengis el mòbil. 📱",
    "Subornant la bàscula virtual... 💸",
    "Això dona per a 3 sèries de sentadilles extres... 🏋️‍♂️",
    "Preguntant-li a un nutricionista italià... 🤌",
    "Separant els greixos bons dels 'regulín'... 🥑",
    "Decidint si això és un 'cheat meal' o un 'cheat day'... 🍕",
  ];

  // State Machine
  const [mode, setMode] = useState<"camera" | "preview" | "loading" | "review">(
    "camera",
  );
  const [imageUri, setImageUri] = useState<string | null>(null);
  const [description, setDescription] = useState("");
  const [facing, setFacing] = useState<"back" | "front">("back");

  // Review State (Editable)
  const [mealName, setMealName] = useState("");
  const [calories, setCalories] = useState("");
  const [protein, setProtein] = useState("");
  const [carbs, setCarbs] = useState("");
  const [fat, setFat] = useState("");
  const [confidence, setConfidence] = useState("");
  const [analysisImagePath, setAnalysisImagePath] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  if (!permission) {
    // Camera permissions are still loading
    return <View style={styles.container} />;
  }

  if (!permission.granted) {
    // Camera permissions are not granted yet
    return (
      <View style={styles.container}>
        <Text style={styles.message}>
          We need your permission to show the camera
        </Text>
        <TouchableOpacity
          onPress={requestPermission}
          style={styles.permissionButton}
        >
          <Text style={styles.permissionButtonText}>Grant Permission</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const takePicture = async () => {
    if (cameraRef.current) {
      try {
        const photo = await cameraRef.current.takePictureAsync({
          quality: 0.7,
          skipProcessing: true,
        });
        setImageUri(photo?.uri || null);
        setMode("preview");
      } catch (error) {
        Alert.alert("Error", "Failed to take picture");
      }
    }
  };

  const getRandomMessage = () => {
    return FUNNY_MESSAGES[Math.floor(Math.random() * FUNNY_MESSAGES.length)];
  };

  const analyzeImage = async () => {
    if (!imageUri) return;

    setMode("loading");
    setLoadingMessage(getRandomMessage());

    try {
      // 1. Submit Job
      const response = await userApi.analyzeImage(imageUri, description);

      // Check if we got 202 Accepted with job_id
      const jobId = response.data.job_id;
      if (!jobId) {
        throw new Error("No Job ID returned from analysis request.");
      }

      // 2. Poll for results
      const pollInterval = setInterval(async () => {
        try {
          // Update message occasionally
          setLoadingMessage(getRandomMessage());

          const jobResponse = await userApi.getAnalysisJob(jobId);
          const jobData = jobResponse.data;

          if (jobData.status === "completed") {
            clearInterval(pollInterval);
            const result = jobData.response;

            if (!result) {
              Alert.alert("Error", "Analysis completed but no result found.");
              setMode("preview");
              return;
            }

            // Populate review state
            setMealName(result.name || "Unknown Meal");
            setCalories(result.calories?.toString() || "0");
            setProtein(result.protein_g?.toString() || "0");
            setCarbs(result.carbs_g?.toString() || "0");
            setFat(result.fat_g?.toString() || "0");
            setConfidence(result.confidence_level || "low");
            setAnalysisImagePath(jobData.image_path || "");

            setMode("review");
          } else if (jobData.status === "failed") {
            clearInterval(pollInterval);
            Alert.alert("Error", jobData.error_message || "Analysis failed.");
            setMode("preview");
          }
          // If 'pending', continue polling
        } catch (pollError) {
          console.error("Polling error:", pollError);
          // Don't verify too aggressively, allow retries?
          // For now, let's stop on error to avoid infinite loops if network is down
          clearInterval(pollInterval);
          Alert.alert("Error", "Failed to check analysis status.");
          setMode("preview");
        }
      }, 3000); // Poll every 3 seconds
    } catch (error: any) {
      console.error("Analysis failed:", error);
      let msg = "Failed to submit analysis.";
      if (error.response?.status === 413) msg = "Image is too large.";
      if (error.response?.data?.error) msg = error.response.data.error;
      Alert.alert("Error", msg);
      setMode("preview");
    }
  };

  const handleSave = async () => {
    if (!mealName.trim()) {
      Alert.alert("Validation Error", "Please enter a meal name.");
      return;
    }

    setIsSaving(true);
    const confidenceMap: Record<string, number> = {
      high: 0.9,
      medium: 0.7,
      low: 0.5,
    };

    try {
      await addMeal({
        name: mealName,
        description: description,
        image_path: analysisImagePath || "--",
        calories: parseInt(calories) || 0,
        protein_g: parseFloat(protein) || 0,
        carbs_g: parseFloat(carbs) || 0,
        fat_g: parseFloat(fat) || 0,
        confidence_level: confidenceMap[confidence?.toLowerCase()] || 0.5,
      });

      // Reset and Navigate to Dashboard
      reset();
      router.replace("/(tabs)");
      // Note: replace to tabs root (Dashboard)
    } catch (error) {
      Alert.alert("Save Failed", "Could not save the meal. Please try again.");
    } finally {
      setIsSaving(false);
    }
  };

  const reset = () => {
    setImageUri(null);
    setDescription("");
    setMealName("");
    setCalories("");
    setProtein("");
    setCarbs("");
    setFat("");
    setAnalysisImagePath("");
    setMode("camera");
  };

  const toggleCamera = () => {
    setFacing((current) => (current === "back" ? "front" : "back"));
  };

  // --- RENDERERS ---

  if (mode === "camera") {
    return (
      <View style={styles.container}>
        <CameraView style={styles.camera} facing={facing} ref={cameraRef}>
          <View style={styles.cameraControls}>
            <TouchableOpacity style={styles.flipButton} onPress={toggleCamera}>
              <Ionicons name="camera-reverse" size={30} color="white" />
            </TouchableOpacity>

            <TouchableOpacity
              style={styles.captureButton}
              onPress={takePicture}
            >
              <View style={styles.captureInner} />
            </TouchableOpacity>

            <View style={styles.placeholder} />
          </View>
        </CameraView>
      </View>
    );
  }

  if (mode === "preview") {
    return (
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={styles.container}
      >
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <Image source={{ uri: imageUri! }} style={styles.previewImage} />

          <Text style={styles.label}>Optional Description:</Text>
          <TextInput
            style={styles.input}
            placeholder="e.g. Added extra olive oil..."
            placeholderTextColor="#666"
            value={description}
            onChangeText={setDescription}
            multiline
          />

          <View style={styles.buttonRow}>
            <TouchableOpacity style={styles.secondaryButton} onPress={reset}>
              <Text style={styles.secondaryButtonText}>Retake</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={styles.primaryButton}
              onPress={analyzeImage}
            >
              <Text style={styles.primaryButtonText}>Analyze Meal</Text>
            </TouchableOpacity>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    );
  }

  if (mode === "loading") {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color="#34C759" />
        <Text style={styles.loadingText}>{loadingMessage}</Text>
        <Text style={styles.loadingSubText}>
          This might take a few seconds...
        </Text>
      </View>
    );
  }

  if (mode === "review") {
    const isHighConfidence = confidence?.toLowerCase() === "high";
    const confidenceColor = isHighConfidence ? "#34C759" : "#FF9500";

    return (
      <SafeAreaView style={styles.container}>
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          style={{ flex: 1 }}
        >
          <ScrollView contentContainerStyle={styles.resultContent}>
            <Image source={{ uri: imageUri! }} style={styles.resultImage} />

            <View style={styles.card}>
              <Text style={styles.cardHeader}>Review & Edit</Text>

              <View style={styles.formGroup}>
                <Text style={styles.label}>Meal Name</Text>
                <TextInput
                  style={styles.inputReview}
                  value={mealName}
                  onChangeText={setMealName}
                  placeholder="Meal Name"
                  placeholderTextColor="#666"
                />
              </View>

              <View style={styles.formGroup}>
                <Text style={styles.label}>Calories (kcal)</Text>
                <TextInput
                  style={[styles.inputReview, styles.highlightInput]}
                  value={calories}
                  onChangeText={setCalories}
                  keyboardType="numeric"
                  placeholder="0"
                  placeholderTextColor="#666"
                />
              </View>

              <View style={styles.macrosGrid}>
                <View style={styles.macroItemEdit}>
                  <Text style={styles.macroLabel}>Protein (g)</Text>
                  <TextInput
                    style={styles.inputMacro}
                    value={protein}
                    onChangeText={setProtein}
                    keyboardType="numeric"
                  />
                </View>
                <View style={styles.macroItemEdit}>
                  <Text style={styles.macroLabel}>Carbs (g)</Text>
                  <TextInput
                    style={styles.inputMacro}
                    value={carbs}
                    onChangeText={setCarbs}
                    keyboardType="numeric"
                  />
                </View>
                <View style={styles.macroItemEdit}>
                  <Text style={styles.macroLabel}>Fat (g)</Text>
                  <TextInput
                    style={styles.inputMacro}
                    value={fat}
                    onChangeText={setFat}
                    keyboardType="numeric"
                  />
                </View>
              </View>

              <Text style={[styles.confidenceText, { color: confidenceColor }]}>
                AI Confidence: {confidence}
              </Text>
            </View>

            <View style={styles.buttonRow}>
              <TouchableOpacity
                style={styles.secondaryButton}
                onPress={() => setMode("preview")}
              >
                <Text style={styles.secondaryButtonText}>Back</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[
                  styles.primaryButton,
                  isSaving && styles.disabledButton,
                ]}
                onPress={handleSave}
                disabled={isSaving}
              >
                {isSaving ? (
                  <ActivityIndicator color="#FFF" />
                ) : (
                  <Text style={styles.primaryButtonText}>Confirm & Log</Text>
                )}
              </TouchableOpacity>
            </View>
          </ScrollView>
        </KeyboardAvoidingView>
      </SafeAreaView>
    );
  }

  return null;
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#121212" },
  centerContainer: {
    flex: 1,
    backgroundColor: "#121212",
    alignItems: "center",
    justifyContent: "center",
  },
  message: {
    textAlign: "center",
    paddingBottom: 10,
    color: "#FFF",
    fontSize: 18,
  },
  camera: { flex: 1 },
  cameraControls: {
    position: "absolute",
    bottom: 50,
    left: 0,
    right: 0,
    flexDirection: "row",
    justifyContent: "space-around",
    alignItems: "center",
  },
  flipButton: { padding: 10 },
  captureButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: "rgba(255,255,255,0.3)",
    alignItems: "center",
    justifyContent: "center",
  },
  captureInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: "#FFF",
  },
  placeholder: { width: 50 },
  permissionButton: {
    backgroundColor: "#007AFF",
    padding: 15,
    borderRadius: 10,
  },
  permissionButtonText: { color: "#FFF", fontWeight: "bold" },
  scrollContent: { padding: 20, alignItems: "center" },
  previewImage: {
    width: "100%",
    height: 300,
    borderRadius: 12,
    marginBottom: 20,
  },
  label: {
    alignSelf: "flex-start",
    color: "#AAA",
    marginBottom: 8,
    fontSize: 14,
  },
  input: {
    width: "100%",
    backgroundColor: "#2C2C2C",
    color: "#FFF",
    padding: 15,
    borderRadius: 10,
    minHeight: 80,
    textAlignVertical: "top",
    marginBottom: 20,
  },
  buttonRow: {
    flexDirection: "row",
    width: "100%",
    justifyContent: "space-between",
    marginBottom: 40,
  },
  primaryButton: {
    backgroundColor: "#34C759",
    padding: 16,
    borderRadius: 10,
    flex: 1,
    marginLeft: 10,
    alignItems: "center",
  },
  primaryButtonText: { color: "#FFF", fontWeight: "bold", fontSize: 16 },
  secondaryButton: {
    backgroundColor: "#3A3A3C",
    padding: 16,
    borderRadius: 10,
    flex: 1,
    marginRight: 10,
    alignItems: "center",
  },
  secondaryButtonText: { color: "#FFF", fontWeight: "600", fontSize: 16 },
  loadingText: {
    color: "#FFF",
    marginTop: 20,
    fontSize: 18,
    fontWeight: "bold",
  },
  loadingSubText: { color: "#AAA", marginTop: 5 },
  resultContent: { padding: 20 },
  resultImage: {
    width: "100%",
    height: 200,
    borderRadius: 16,
    marginBottom: -20,
    zIndex: 1,
  },
  card: {
    backgroundColor: "#1E1E1E",
    borderRadius: 20,
    padding: 24,
    paddingTop: 30,
    marginBottom: 20,
    zIndex: 0,
  },
  cardHeader: {
    color: "#FFF",
    fontSize: 22,
    fontWeight: "bold",
    marginBottom: 20,
    textAlign: "center",
  },
  formGroup: { marginBottom: 16 },
  inputReview: {
    backgroundColor: "#2C2C2C",
    color: "#FFF",
    padding: 12,
    borderRadius: 8,
    fontSize: 16,
  },
  highlightInput: {
    borderColor: "#34C759",
    borderWidth: 1,
    fontWeight: "bold",
    fontSize: 18,
    textAlign: "center",
  },
  macrosGrid: {
    flexDirection: "row",
    justifyContent: "space-between",
    width: "100%",
    marginBottom: 16,
  },
  macroItemEdit: { flex: 1, marginHorizontal: 4 },
  macroLabel: {
    color: "#AAA",
    fontSize: 12,
    marginBottom: 4,
    textAlign: "center",
  },
  inputMacro: {
    backgroundColor: "#2C2C2C",
    color: "#FFF",
    padding: 10,
    borderRadius: 8,
    textAlign: "center",
    fontSize: 16,
  },
  confidenceText: {
    textAlign: "center",
    fontSize: 12,
    marginTop: 8,
    fontStyle: "italic",
  },
  disabledButton: { opacity: 0.6 },
});
