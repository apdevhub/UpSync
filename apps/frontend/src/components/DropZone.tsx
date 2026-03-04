import { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { formatBytes, getFileIcon } from '@/utils';

const MAX_SIZE = 50 * 1024 * 1024; // 50 MB

interface Props {
    file: File | null;
    onFile: (f: File) => void;
    onRemove: () => void;
}

export default function DropZone({ file, onFile, onRemove }: Props) {
    const [sizeError, setSizeError] = useState<string | null>(null);

    const onDrop = useCallback(
        (accepted: File[], rejected: { errors: { code: string; message: string }[] }[]) => {
            setSizeError(null);
            if (rejected.length > 0) {
                const err = rejected[0]?.errors[0];
                if (err?.code === 'file-too-large') {
                    setSizeError('File exceeds the 50 MB limit. Please choose a smaller file.');
                } else {
                    setSizeError(err?.message ?? 'Invalid file.');
                }
                return;
            }
            if (accepted[0]) onFile(accepted[0]);
        },
        [onFile],
    );

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        maxSize: MAX_SIZE,
        multiple: false,
    });

    if (file) {
        return (
            <div className="file-preview" role="region" aria-label="Selected file">
                <div className="file-type-icon" aria-hidden="true">
                    {getFileIcon(file.type)}
                </div>
                <div className="file-info">
                    <div className="file-name" title={file.name}>{file.name}</div>
                    <div className="file-size">{formatBytes(file.size)}</div>
                </div>
                <button
                    className="remove-btn"
                    onClick={onRemove}
                    aria-label="Remove selected file"
                    type="button"
                >
                    ✕
                </button>
            </div>
        );
    }

    return (
        <>
            <div
                {...getRootProps()}
                id="drop-zone"
                className={`upload-zone${isDragActive ? ' drag-active' : ''}`}
                aria-label="File drop zone"
            >
                <input {...getInputProps()} id="file-input" aria-label="File input" />
                <div className="upload-icon-wrap" aria-hidden="true">
                    {isDragActive ? '📂' : '☁️'}
                </div>
                <h2>{isDragActive ? "Drop it like it's hot!" : 'Drop your file here'}</h2>
                <p>
                    or <span className="browse">browse to choose</span> a file
                </p>
                <div className="upload-limit">Maximum file size: 50 MB</div>
            </div>
            {sizeError && (
                <p role="alert" style={{ color: 'var(--danger)', fontSize: 13, marginTop: 10, textAlign: 'center' }}>
                    ⚠ {sizeError}
                </p>
            )}
        </>
    );
}
