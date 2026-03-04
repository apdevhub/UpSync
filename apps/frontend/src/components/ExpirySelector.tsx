import type { ExpiryOption } from '@/types';

const EXPIRY_OPTIONS: ExpiryOption[] = [
    { label: '1 Hour', value: '1h' },
    { label: '6 Hours', value: '6h' },
    { label: '12 Hours', value: '12h' },
    { label: '1 Day', value: '24h' },
    { label: '3 Days', value: '3d' },
    { label: '7 Days', value: '7d' },
];

interface Props {
    value: string;
    onChange: (v: string) => void;
}

export default function ExpirySelector({ value, onChange }: Props) {
    return (
        <div className="expiry-row" role="group" aria-label="File expiry duration">
            <span className="expiry-label">⏱ Expires in</span>
            <div className="expiry-options">
                {EXPIRY_OPTIONS.map((opt) => (
                    <button
                        key={opt.value}
                        type="button"
                        id={`expiry-${opt.value}`}
                        className={`expiry-btn${value === opt.value ? ' active' : ''}`}
                        onClick={() => onChange(opt.value)}
                        aria-pressed={value === opt.value}
                    >
                        {opt.label}
                    </button>
                ))}
            </div>
        </div>
    );
}
