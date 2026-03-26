<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from './lib/api';
  import type { User, Category, Bookmark, AuthResponse } from './lib/types';
  import Login from './components/Login.svelte';
  import Register from './components/Register.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import BookmarkForm from './components/BookmarkForm.svelte';
  import BookmarkList from './components/BookmarkList.svelte';
  import ChangePassword from './components/ChangePassword.svelte';

  let user = $state<User | null>(null);
  let view = $state<'login' | 'register' | 'dashboard'>('login');
  let categories = $state<Category[]>([]);
  let urls = $state<Bookmark[]>([]);
  let selectedCategory = $state<Category | null>(null);
  let searchQuery = $state('');
  let editingUrl = $state<Bookmark | null>(null);
  let loading = $state(true);
  let error = $state('');
  let showChangePassword = $state(false);

  onMount(async () => {
    const token = api.getToken();
    if (token) {
      try {
        user = await api.getProfile();
        view = 'dashboard';
        await loadData();
      } catch (err: any) {
        api.setToken(null);
        view = 'login';
      }
    } else {
      view = 'login';
    }
    loading = false;
  });

  async function loadData() {
    try {
      const [cats, urlList] = await Promise.all([
        api.listCategories(),
        api.listURLs(selectedCategory?.id, searchQuery || undefined),
      ]);
      categories = cats;
      urls = urlList;
    } catch (err: any) {
      error = 'Failed to load data';
    }
  }

  function handleLogin(data: AuthResponse) {
    user = data.user ?? data;
    view = 'dashboard';
    loadData();
  }

  function handleRegister(data: AuthResponse) {
    user = data.user ?? data;
    view = 'dashboard';
    loadData();
  }

  async function handleLogout() {
    api.setToken(null);
    user = null;
    view = 'login';
    categories = [];
    urls = [];
  }

  async function handleCategorySelect(category: Category | null) {
    selectedCategory = category;
    await loadData();
  }

  async function handleSearch(query: string) {
    searchQuery = query;
    await loadData();
  }

  function handleBookmarkCreate(bookmark: Bookmark) {
    urls = [bookmark, ...urls];
  }

  function handleBookmarkUpdate(updated: Bookmark) {
    urls = urls.map(u => u.id === updated.id ? updated : u);
    editingUrl = null;
  }

  function handleBookmarkDelete(bookmark: Bookmark) {
    urls = urls.filter(u => u.id !== bookmark.id);
  }

  function handleBookmarkEdit(bookmark: Bookmark) {
    editingUrl = bookmark;
  }

  function cancelEdit() {
    editingUrl = null;
  }

  function userInitials(u: User): string {
    return u.name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase();
  }
</script>

{#if loading}
  <div class="min-h-screen flex items-center justify-center bg-base-200">
    <div class="flex flex-col items-center gap-3">
      <div class="w-10 h-10 rounded-xl bg-green-700 flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
      </div>
      <span class="loading loading-dots loading-sm text-green-700"></span>
    </div>
  </div>
{:else if view === 'login'}
  <Login
    onlogin={handleLogin}
    onregister={() => view = 'register'}
  />
{:else if view === 'register'}
  <Register
    onregister={handleRegister}
    onlogin={() => view = 'login'}
  />
{:else if view === 'dashboard'}
  <div class="flex min-h-screen bg-base-200">
    <Sidebar
      {categories}
      {selectedCategory}
      {searchQuery}
      onselect={handleCategorySelect}
      onsearch={handleSearch}
    />

    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top bar -->
      <header class="bg-base-100 border-b border-base-200 px-6 py-3 flex items-center justify-between shrink-0">
        <div>
          <h1 class="font-semibold text-base-content">
            {selectedCategory ? selectedCategory.name : 'All Bookmarks'}
          </h1>
          <p class="text-xs text-base-content/50">{urls.length} bookmark{urls.length !== 1 ? 's' : ''}</p>
        </div>

        <div class="flex items-center gap-3">
          {#if user}
            <button
              class="flex items-center gap-2.5 rounded-lg px-2 py-1 hover:bg-base-200 transition-colors"
              onclick={() => showChangePassword = true}
              title="Change password"
            >
              <div class="w-8 h-8 rounded-full bg-green-700 flex items-center justify-center text-white text-xs font-semibold select-none">
                {userInitials(user)}
              </div>
              <span class="text-sm text-base-content/70 hidden sm:block">{user.name}</span>
            </button>
          {/if}
          <button class="btn btn-ghost btn-sm gap-1.5 text-base-content/60 hover:text-error" onclick={handleLogout}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Sign out
          </button>
        </div>
      </header>

      <!-- Main content -->
      <main class="flex-1 p-6 overflow-y-auto">
        {#if error}
          <div class="alert alert-error mb-5 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{error}</span>
          </div>
        {/if}

        <BookmarkForm
          {categories}
          {editingUrl}
          oncreate={handleBookmarkCreate}
          onupdate={handleBookmarkUpdate}
          oncancel={cancelEdit}
        />

        <BookmarkList
          {urls}
          {categories}
          onedit={handleBookmarkEdit}
          ondelete={handleBookmarkDelete}
        />
      </main>
    </div>
  </div>

  {#if showChangePassword}
    <ChangePassword onclose={() => showChangePassword = false} />
  {/if}
{/if}
