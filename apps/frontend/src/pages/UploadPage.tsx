import { useState, useCallback } from 'react';
import { uploadFile } from '@/api';
import { useToast } from '@/context/ToastContext';
import DropZone from '@/components/DropZone';
import ExpirySelector from '@/components/ExpirySelector';
import { formatBytes, formatDate } from '@/utils';
import type { UploadResponse } from '@/types';

export default function UploadPage() {
    const { showToast } = useToast();
    const [file, setFile] = useState<File | null>(null);
    const [expiresIn, setExpiresIn] = useState<string>('24h');
    const [uploading, setUploading] = useState<boolean>(false);
    const [progress, setProgress] = useState<number>(0);
    const [result, setResult] = useState<UploadResponse | null>(null);
    const [copied, setCopied] = useState<boolean>(false);

    const handleUpload = useCallback(async () => {
        if (!file) return;
        setUploading(true);
        setProgress(0);
        try {
            const { data } = await uploadFile(file, expiresIn, setProgress);
            setResult(data);
            setFile(null);
            showToast('File uploaded! Share the link below.', 'success');
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Upload failed.', 'error');
        } finally {
            setUploading(false);
        }
    }, [file, expiresIn, showToast]);

    const handleCopy = useCallback(async () => {
        if (!result) return;
        try {
            await navigator.clipboard.writeText(result.shareUrl);
            setCopied(true);
            showToast('Link copied to clipboard!', 'success');
            setTimeout(() => setCopied(false), 2500);
        } catch {
            showToast('Could not copy — please copy manually.', 'error');
        }
    }, [result, showToast]);

    const handleReset = () => {
        setResult(null);
        setFile(null);
        setProgress(0);
    };

    return (
        <main className="page" aria-label="Upload page">
            <div className="container">

                {/* ── Hero ─────────────────────────────────────────── */}
                <section className="hero" aria-labelledby="hero-heading">
                    <div className="hero-badge" aria-hidden="true">
                        <span>⚡</span> Fast · Secure · Temporary
                    </div>
                    <h1 id="hero-heading">Share Files<br />Instantly</h1>
                    <p>
                        Upload any file up to 50 MB, get a shareable link, and let it
                        auto-delete on your schedule.
                    </p>
                </section>

                {/* ── Upload Form or Success Card ───────────────────── */}
                {!result ? (
                    <section aria-label="Upload form">
                        <DropZone file={file} onFile={setFile} onRemove={() => setFile(null)} />

                        <ExpirySelector value={expiresIn} onChange={setExpiresIn} />

                        {uploading && (
                            <div className="progress-wrap" role="progressbar" aria-valuenow={progress} aria-valuemin={0} aria-valuemax={100}>
                                <div className="progress-bar-bg">
                                    <div className="progress-bar-fill" style={{ width: `${progress}%` }} />
                                </div>
                                <div className="progress-label">
                                    <span>Uploading to secure storage…</span>
                                    <span>{progress}%</span>
                                </div>
                            </div>
                        )}

                        <button
                            id="upload-btn"
                            className="btn-upload"
                            onClick={handleUpload}
                            disabled={!file || uploading}
                            aria-busy={uploading}
                            type="button"
                        >
                            <span>{uploading ? '⏳ Uploading…' : '🚀 Upload & Get Link'}</span>
                        </button>
                    </section>
                ) : (
                    <section className="success-card" aria-label="Upload success">
                        <div className="success-icon" aria-hidden="true">✅</div>
                        <h2>Upload Complete!</h2>
                        <p>Your file is ready to share. The link will expire automatically.</p>

                        <div className="share-url-box">
                            <div className="share-url-text" aria-label="Share link">{result.shareUrl}</div>
                            <button
                                id="copy-btn"
                                className={`copy-btn${copied ? ' copied' : ''}`}
                                onClick={handleCopy}
                                aria-label="Copy share link"
                                type="button"
                            >
                                {copied ? '✓ Copied!' : '📋 Copy'}
                            </button>
                        </div>

                        <div className="meta-pills" aria-label="File metadata">
                            <div className="meta-pill">📦 {formatBytes(result.size)}</div>
                            <div className="meta-pill">⏱ Expires {formatDate(result.expiresAt)}</div>
                        </div>

                        <button
                            id="upload-another-btn"
                            className="btn-secondary"
                            onClick={handleReset}
                            type="button"
                        >
                            ↑ Upload Another File
                        </button>
                    </section>
                )}
            </div>
        </main>
    );
}
