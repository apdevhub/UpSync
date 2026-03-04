import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getFileMeta, getDownloadUrl } from '@/api';
import { useToast } from '@/context/ToastContext';
import Spinner from '@/components/ui/Spinner';
import { formatBytes, getFileIcon, formatExpiry, formatDate } from '@/utils';
import type { FileMeta } from '@/types';

type PageState = 'loading' | 'ready' | 'expired' | 'not_found' | 'error';

export default function SharePage() {
    const { id } = useParams<{ id: string }>();
    const { showToast } = useToast();
    const [state, setState] = useState<PageState>('loading');
    const [meta, setMeta] = useState<FileMeta | null>(null);
    const [downloading, setDownloading] = useState<boolean>(false);

    useEffect(() => {
        if (!id) { setState('not_found'); return; }

        getFileMeta(id)
            .then(({ data }) => {
                setMeta(data);
                setState('ready');
            })
            .catch((err: Error) => {
                // Axios interceptor converts HTTP 410 → Error with specific message
                if (err.message.toLowerCase().includes('expired') || err.message.includes('deleted')) {
                    setState('expired');
                } else if (err.message.toLowerCase().includes('not found')) {
                    setState('not_found');
                } else {
                    setState('error');
                }
            });
    }, [id]);

    const handleDownload = async () => {
        if (!id) return;
        setDownloading(true);
        try {
            const { data } = await getDownloadUrl(id);
            const a = document.createElement('a');
            a.href = data.downloadUrl;
            a.download = data.fileName;
            document.body.appendChild(a);
            a.click();
            a.remove();
            showToast('Download started!', 'success');
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Download failed.', 'error');
        } finally {
            setDownloading(false);
        }
    };

    /* ── Loading ─────────────────────────────────────────────── */
    if (state === 'loading') {
        return (
            <div className="share-page">
                <div className="status-card">
                    <Spinner />
                    <p style={{ color: 'var(--text-secondary)' }}>Loading file info…</p>
                </div>
            </div>
        );
    }

    /* ── Expired ─────────────────────────────────────────────── */
    if (state === 'expired') {
        return (
            <div className="share-page">
                <div className="status-card">
                    <span className="status-icon">⏰</span>
                    <h2>Link Expired</h2>
                    <p>This share link has passed its expiry date.<br />The file has been automatically deleted.</p>
                    <Link to="/" className="go-home">🏠 Go to UpSync</Link>
                </div>
            </div>
        );
    }

    /* ── Not Found / Error ───────────────────────────────────── */
    if (state === 'not_found' || state === 'error') {
        return (
            <div className="share-page">
                <div className="status-card">
                    <span className="status-icon">🔍</span>
                    <h2>{state === 'not_found' ? 'File Not Found' : 'Something Went Wrong'}</h2>
                    <p>
                        {state === 'not_found'
                            ? "This link doesn't exist or the file may have been deleted."
                            : 'An unexpected error occurred. Please try again later.'}
                    </p>
                    <Link to="/" className="go-home">🏠 Go to UpSync</Link>
                </div>
            </div>
        );
    }

    /* ── Ready ───────────────────────────────────────────────── */
    if (!meta) return null;

    return (
        <div className="share-page">
            <div className="bg-orbs" aria-hidden="true" />

            <article className="file-card" aria-label={`Download ${meta.originalName}`}>
                {/* Brand link */}
                <Link
                    to="/"
                    style={{ textDecoration: 'none', display: 'inline-flex', alignItems: 'center', gap: 8, marginBottom: 28, opacity: 0.7 }}
                    aria-label="Go to UpSync homepage"
                >
                    <div style={{ width: 24, height: 24, background: 'linear-gradient(135deg, #7c3aed, #06b6d4)', borderRadius: 6, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                            <path d="M12 4L12 16M12 4L8 8M12 4L16 8" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
                            <path d="M5 16C3.34 16 2 17.34 2 19C2 20.66 3.34 22 5 22H19C20.66 22 22 20.66 22 19C22 17.34 20.66 16 19 16" stroke="white" strokeWidth="2" strokeLinecap="round" />
                        </svg>
                    </div>
                    <span style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-secondary)' }}>UpSync</span>
                </Link>

                <div className="file-card-icon" aria-hidden="true">
                    {getFileIcon(meta.mimeType)}
                </div>

                <h1 id="file-name">{meta.originalName}</h1>
                <div className="sub">Shared via UpSync</div>

                <div className="stats-grid" aria-label="File details">
                    <div className="stat-box">
                        <div className="stat-label">File Size</div>
                        <div className="stat-value">{formatBytes(meta.size)}</div>
                    </div>
                    <div className="stat-box">
                        <div className="stat-label">Type</div>
                        <div className="stat-value" style={{ fontSize: 13 }}>{meta.mimeType}</div>
                    </div>
                    <div className="stat-box">
                        <div className="stat-label">Uploaded</div>
                        <div className="stat-value" style={{ fontSize: 13 }}>{formatDate(meta.createdAt)}</div>
                    </div>
                    <div className="stat-box">
                        <div className="stat-label">Status</div>
                        <div className="stat-value" style={{ color: 'var(--success)' }}>● Active</div>
                    </div>
                </div>

                <button
                    id="download-btn"
                    className="btn-download"
                    onClick={handleDownload}
                    disabled={downloading}
                    aria-busy={downloading}
                    type="button"
                >
                    {downloading ? '⏳ Preparing…' : '⬇ Download File'}
                </button>

                <div className="expiry-notice" aria-live="polite">
                    ⚠ {formatExpiry(meta.expiresAt)}
                </div>
            </article>
        </div>
    );
}
