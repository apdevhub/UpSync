import axios, { AxiosError } from 'axios';
import type { FileMeta, UploadResponse, DownloadResponse, ApiError } from '@/types';

// ─── Axios Instance ───────────────────────────────────────────
const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080',
    timeout: 60_000,
});

// ─── Response Interceptor — normalise errors ──────────────────
api.interceptors.response.use(
    (res) => res,
    (err: AxiosError<ApiError>) => {
        const message =
            err.response?.data?.error ??
            (err.code === 'ECONNABORTED' ? 'Request timed out. Please try again.' : 'Network error. Is the server running?');
        return Promise.reject(new Error(message));
    },
);

// ─── API Methods ──────────────────────────────────────────────
export const uploadFile = (
    file: File,
    expiresIn: string,
    onProgress?: (pct: number) => void,
) => {
    const fd = new FormData();
    fd.append('file', file);
    fd.append('expiresIn', expiresIn);

    return api.post<UploadResponse>('/api/files/upload', fd, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (e) => {
            if (onProgress && e.total) {
                onProgress(Math.round((e.loaded / e.total) * 100));
            }
        },
    });
};

export const getFileMeta = (id: string) =>
    api.get<FileMeta>(`/api/files/${id}`);

export const getDownloadUrl = (id: string) =>
    api.get<DownloadResponse>(`/api/files/${id}/download`);

export default api;
