export interface User {
  id: number;
  username: string;
  name: string;
}

export interface Category {
  id: number;
  name: string;
  description?: string;
  color?: string;
}

export interface Bookmark {
  id: number;
  url: string;
  title?: string;
  description?: string;
  is_pinned: boolean;
  category_id?: number;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  user: User;
}
