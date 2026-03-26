<script lang="ts">
  import { api } from '../lib/api';
  import type { Bookmark, Category } from '../lib/types';

  let { urls = [], categories = [], onedit, ondelete }: {
    urls?: Bookmark[];
    categories?: Category[];
    onedit?: (bookmark: Bookmark) => void;
    ondelete?: (bookmark: Bookmark) => void;
  } = $props();

  function getCategoryName(categoryId: number): string {
    const category = categories.find(c => c.id === categoryId);
    return category?.name ?? 'Unknown';
  }

  function getDomain(url: string): string {
    try {
      return new URL(url).hostname.replace('www.', '');
    } catch {
      return url;
    }
  }

  function getDomainInitial(url: string): string {
    return getDomain(url).charAt(0).toUpperCase();
  }

  // Deterministic color from domain string
  const AVATAR_COLORS = [
    'bg-violet-500', 'bg-blue-500', 'bg-emerald-500', 'bg-orange-500',
    'bg-rose-500', 'bg-teal-500', 'bg-indigo-500', 'bg-amber-500',
  ];

  function getAvatarColor(url: string): string {
    const domain = getDomain(url);
    let hash = 0;
    for (let i = 0; i < domain.length; i++) hash = domain.charCodeAt(i) + ((hash << 5) - hash);
    return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length];
  }

  async function deleteURL(url: Bookmark) {
    if (!confirm(`Delete bookmark "${url.url}"?`)) return;

    try {
      await api.deleteURL(url.id);
      ondelete?.(url);
    } catch (err: any) {
      alert(err.message);
    }
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }
</script>

{#if urls.length === 0}
  <div class="flex flex-col items-center justify-center py-24 text-base-content/40">
    <div class="w-16 h-16 rounded-2xl bg-base-200 flex items-center justify-center mb-4">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
      </svg>
    </div>
    <p class="font-medium text-base-content/60 mb-1">No bookmarks yet</p>
    <p class="text-sm">Add your first bookmark above to get started</p>
  </div>
{:else}
  <div class="flex flex-col gap-2">
    {#each urls as url}
      <div class="group flex items-center gap-4 px-4 py-3.5 bg-base-100 rounded-xl border border-base-200 hover:border-violet-200 hover:shadow-sm transition-all">
        <!-- Domain avatar -->
        <div class="shrink-0 w-9 h-9 rounded-lg {getAvatarColor(url.url)} flex items-center justify-center text-white font-semibold text-sm select-none">
          {getDomainInitial(url.url)}
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <div class="flex items-baseline gap-2">
            <a
              href={url.url}
              target="_blank"
              rel="noopener noreferrer"
              class="font-medium text-base-content hover:text-violet-600 truncate transition-colors"
            >
              {getDomain(url.url)}
            </a>
            <span class="text-xs text-base-content/40 truncate hidden sm:block">{url.url}</span>
          </div>
          {#if url.description}
            <p class="text-sm text-base-content/60 truncate mt-0.5">{url.description}</p>
          {/if}
        </div>

        <!-- Meta -->
        <div class="hidden md:flex items-center gap-3 shrink-0">
          <span class="badge badge-ghost badge-sm text-xs">{getCategoryName(url.category_id)}</span>
          <span class="text-xs text-base-content/40">{formatDate(url.created_at)}</span>
        </div>

        <!-- Actions -->
        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
          <button
            class="btn btn-ghost btn-sm btn-circle text-base-content/50 hover:text-violet-600"
            onclick={() => onedit?.(url)}
            title="Edit"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>
          <button
            class="btn btn-ghost btn-sm btn-circle text-base-content/50 hover:text-error"
            onclick={() => deleteURL(url)}
            title="Delete"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>
    {/each}
  </div>
{/if}
