interface Props {
    className?: string;
}

export default function Spinner({ className = '' }: Props) {
    return <div className={`spinner ${className}`} role="status" aria-label="Loading" />;
}
