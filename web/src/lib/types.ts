// API Response types
export interface APIResponse<T = any> {
  code: string;
  message: string;
  data?: T;
}

// User types
export interface User {
  id: string;
  email: string;
  username: string;
  first_name?: string | null;
  last_name?: string | null;
  created_at: string;
  updated_at: string;
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  first_name?: string | null;
  last_name?: string | null;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

// Mental Health Record types
export interface MentalHealthRecord {
  id: string;
  user_id: string;
  happy_level: number;
  energy_level: number;
  notes?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CreateMentalHealthRecordRequest {
  happy_level: number;
  energy_level: number;
  notes?: string;
  status: string;
}

export interface UpdateMentalHealthRecordRequest {
  happy_level?: number;
  energy_level?: number;
  notes?: string;
  status: string;
}

// Heatmap types
export interface HeatmapDataPoint {
  happy_level: number;
  energy_level: number;
  count: number;
}

export interface DateRange {
  started_at?: string;
  ended_at?: string;
}

export interface MentalHealthHeatmapResponse {
  data: Record<string, HeatmapDataPoint>;
  total_records: number;
  date_range: DateRange;
}

// Quote types
export interface Quote {
  id: string;
  content: string;
  author: string;
  created_at: string;
  updated_at: string;
}
