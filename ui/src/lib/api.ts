const API_URL = import.meta.env.PROD ? '' : ''; // Empty for same origin in both dev and prod

class ApiService {
  private token: string | null = null;

  setToken(token: string | null) {
    this.token = token;
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }

  getToken(): string | null {
    if (!this.token) {
      this.token = localStorage.getItem('token');
    }
    return this.token;
  }

  private async request(endpoint: string, options: RequestInit = {}): Promise<any> {
    const url = `${API_URL}/api${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>,
    };

    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (response.status === 204) {
      return null;
    }

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'An error occurred');
    }

    return data;
  }

  // Auth
  async register(username: string, password: string, name: string) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password, name }),
    });
  }

  async login(username: string, password: string) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
  }

  async logout(refreshToken: string) {
    return this.request('/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  // Profile
  async getProfile() {
    return this.request('/profile');
  }

  // Categories
  async listCategories() {
    return this.request('/category');
  }

  async getCategory(id: number) {
    return this.request(`/category/${id}`);
  }

  async createCategory(name: string, description?: string) {
    return this.request('/category', {
      method: 'POST',
      body: JSON.stringify({ name, description }),
    });
  }

  async updateCategory(id: number, name: string, description?: string) {
    return this.request(`/category/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name, description }),
    });
  }

  async deleteCategory(id: number) {
    return this.request(`/category/${id}`, {
      method: 'DELETE',
    });
  }

  // URLs
  async listURLs(categoryId?: number, search?: string) {
    let query = '';
    if (categoryId) query += `?category_id=${categoryId}`;
    if (search) query += `${query ? '&' : '?'}search=${encodeURIComponent(search)}`;
    return this.request(`/url${query}`);
  }

  async getURL(id: number) {
    return this.request(`/url/${id}`);
  }

  async createURL(url: string, categoryId: number, description?: string) {
    return this.request('/url', {
      method: 'POST',
      body: JSON.stringify({ url, category_id: categoryId, description }),
    });
  }

  async updateURL(id: number, url: string, categoryId: number, description?: string) {
    return this.request(`/url/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ url, category_id: categoryId, description }),
    });
  }

  async deleteURL(id: number) {
    return this.request(`/url/${id}`, {
      method: 'DELETE',
    });
  }
}

export const api = new ApiService();
