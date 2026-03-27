<script lang="ts">
  import { api } from '../lib/api';
  import type { Bookmark, Category } from '../lib/types';

  let { urls = [], categories = [], onedit, ondelete }: {
    urls?: Bookmark[];
    categories?: Category[];
    onedit?: (bookmark: Bookmark) => void;
    ondelete?: (bookmark: Bookmark) => void;
  } = $props();

  function getCategoryName(categoryId?: number): string {
    if (!categoryId) return '';
    const category = categories.find(c => c.id === categoryId);
    return category?.name ?? '';
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

  // Deterministic color from domain string - dark green palette
  const AVATAR_COLORS = [
    'bg-green-700', 'bg-green-800', 'bg-emerald-700', 'bg-emerald-800',
    'bg-teal-700', 'bg-teal-800', 'bg-green-900', 'bg-emerald-900',
    'bg-lime-700', 'bg-lime-800', 'bg-green-600', 'bg-teal-900',
  ];

  function getAvatarColor(url: string): string {
    const domain = getDomain(url);
    let hash = 0;
    for (let i = 0; i < domain.length; i++) hash = domain.charCodeAt(i) + ((hash << 5) - hash);
    return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length];
  }

  async function deleteURL(url: Bookmark) {
    if (!confirm(`Delete bookmark "${url.title || getDomain(url.url)}"?`)) return;

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
  <div class="grid grid-cols-3 gap-4 items-start">
    {#each urls as url (url.id)}
      <div class="group bg-base-100 rounded-xl border {url.is_pinned ? 'border-green-300 shadow-sm' : 'border-base-200'} hover:border-green-300 hover:shadow-lg transition-all flex flex-col">
        <!-- Card Header with Avatar -->
        <div class="p-4 pb-3">
          <div class="flex items-start gap-3">
            <!-- Domain avatar -->
            <div class="shrink-0 w-10 h-10 rounded-lg {getAvatarColor(url.url)} flex items-center justify-center text-white font-semibold text-sm select-none shadow-sm">
              {getDomainInitial(url.url)}
            </div>

            <!-- Title & Domain -->
            <div class="flex-1 min-w-0 pt-0.5">
              <a
                href={url.url}
                target="_blank"
                rel="noopener noreferrer"
                class="font-semibold text-base-content hover:text-green-700 transition-colors line-clamp-1"
              >
                {url.title || getDomain(url.url)}
              </a>
              <p class="text-xs text-base-content/40 line-clamp-1 mt-0.5">{url.url}</p>
            </div>

            <!-- Pin indicator + Actions -->
            <div class="flex items-center gap-0.5 shrink-0 -mr-1 -mt-1">
              {#if url.is_pinned}
                <span class="text-green-700 opacity-60" title="Pinned">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M16 12V4h1V2H7v2h1v8l-2 2v2h5.2v6h1.6v-6H18v-2l-2-2z"/>
                  </svg>
                </span>
              {/if}
              <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                <button
                  class="btn btn-ghost btn-sm btn-circle h-8 w-8 min-h-8 text-base-content/50 hover:text-green-700"
                  onclick={() => onedit?.(url)}
                  title="Edit"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button
                  class="btn btn-ghost btn-sm btn-circle h-8 w-8 min-h-8 text-base-content/50 hover:text-error"
                  onclick={() => deleteURL(url)}
                  title="Delete"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        {#if url.description}
          <div class="px-4 pb-3 flex-1">
            <p class="text-sm text-base-content/70 leading-relaxed">{url.description}</p>
          </div>
        {/if}

        <!-- Card Footer -->
        <div class="px-4 py-3 border-t border-base-200 bg-base-50/50 rounded-b-xl">
          <div class="flex items-center justify-between">
            {#if url.category_id}
              <span class="badge badge-sm text-xs font-medium bg-green-800 text-white border-0">
                {getCategoryName(url.category_id)}
              </span>
            {:else}
              <span class="text-xs text-base-content/30 italic">Uncategorized</span>
            {/if}
            <span class="text-xs text-base-content/40">{formatDate(url.created_at)}</span>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
