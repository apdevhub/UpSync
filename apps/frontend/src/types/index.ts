// ─── File Metadata from API ──────────────────────────────────
export interface FileMeta {
    id: string;
    originalName: string;
    mimeType: string;
    size: number;
    expiresAt: string;
    createdAt: string;
}

// ─── Upload Response ─────────────────────────────────────────
export interface UploadResponse {
    id: string;
    originalName: string;
    size: number;
    mimeType: string;
    expiresAt: string;
    shareUrl: string;
}

// ─── Download URL Response ───────────────────────────────────
export interface DownloadResponse {
    downloadUrl: string;
    fileName: string;
}

// ─── API Error Shape ─────────────────────────────────────────
export interface ApiError {
    error: string;
    code?: string;
}

// ─── Expiry Option ───────────────────────────────────────────
export interface ExpiryOption {
    label: string;
    value: string;
}

// ─── Toast Notification ──────────────────────────────────────
export type ToastType = 'success' | 'error' | 'info';

export interface ToastItem {
    id: string;
    message: string;
    type: ToastType;
}
