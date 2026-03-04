import { useToast } from '@/context/ToastContext';
import type { ToastItem } from '@/types';

function ToastItem({ toast, onClose }: { toast: ToastItem; onClose: () => void }) {
    const icons: Record<string, string> = { success: '✅', error: '❌', info: 'ℹ️' };
    return (
        <div className={`toast ${toast.type}`} role="alert">
            <span>{icons[toast.type]}</span>
            <span>{toast.message}</span>
            <button
                onClick={onClose}
                style={{ marginLeft: 12, background: 'none', border: 'none', color: 'inherit', cursor: 'pointer', opacity: 0.6 }}
                aria-label="Dismiss"
            >
                ✕
            </button>
        </div>
    );
}

export default function ToastContainer() {
    const { toasts, removeToast } = useToast();
    if (toasts.length === 0) return null;

    return (
        <div className="toast-wrap" role="region" aria-label="Notifications">
            {toasts.map((t) => (
                <ToastItem key={t.id} toast={t} onClose={() => removeToast(t.id)} />
            ))}
        </div>
    );
}
