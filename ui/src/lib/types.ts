export interface User {
  id: number;
  username: string;
  name: string;
}

export interface Category {
  id: number;
  name: string;
  description?: string;
}

export interface Bookmark {
  id: number;
  url: string;
  description?: string;
  category_id: number;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  user: User;
}
