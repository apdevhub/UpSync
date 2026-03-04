// ─── Format bytes to human-readable string ───────────────────
export function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'] as const;
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
}

// ─── Get emoji icon from MIME type ───────────────────────────
export function getFileIcon(mimeType: string): string {
    if (mimeType.startsWith('image/')) return '🖼️';
    if (mimeType.startsWith('video/')) return '🎬';
    if (mimeType.startsWith('audio/')) return '🎵';
    if (mimeType.includes('pdf')) return '📄';
    if (/zip|rar|7z|tar|gz/.test(mimeType)) return '📦';
    if (/word|document/.test(mimeType)) return '📝';
    if (/sheet|excel/.test(mimeType)) return '📊';
    if (/presentation|powerpoint/.test(mimeType)) return '📑';
    if (mimeType.startsWith('text/')) return '📃';
    return '📁';
}

// ─── Format expiry as relative or absolute ───────────────────
export function formatExpiry(isoDate: string): string {
    const d = new Date(isoDate);
    const diffMs = d.getTime() - Date.now();
    if (diffMs <= 0) return 'Expired';
    const diffH = Math.floor(diffMs / 3_600_000);
    const diffM = Math.floor((diffMs % 3_600_000) / 60_000);
    if (diffH >= 24) {
        const days = Math.floor(diffH / 24);
        return `${days} day${days > 1 ? 's' : ''} remaining`;
    }
    if (diffH > 0) return `${diffH}h ${diffM}m remaining`;
    return `${diffM} min remaining`;
}

// ─── Format ISO date to locale string ────────────────────────
export function formatDate(isoDate: string): string {
    return new Date(isoDate).toLocaleString(undefined, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
}

// ─── Generate a unique ID ─────────────────────────────────────
export function uid(): string {
    return Math.random().toString(36).slice(2, 9);
}
