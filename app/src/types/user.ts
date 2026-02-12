export interface Macros {
  id: string;
  user_id: string;
  timestamp: string;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  active: boolean;
}

export interface CreateMacrosRequest {
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
}
