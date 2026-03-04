import { Link, useLocation } from 'react-router-dom';

export default function Navbar() {
    const { pathname } = useLocation();
    const isShare = pathname.startsWith('/share/');

    return (
        <nav className="navbar" aria-label="Main navigation">
            <div className="navbar-inner">
                <Link to="/" className="navbar-logo" id="nav-logo" aria-label="UpSync home">
                    <div className="logo-icon" aria-hidden="true">
                        <img src="/upsync-logo.svg?v=3" alt="UpSync Icon" width="28" height="28" style={{ display: 'block', borderRadius: '4px' }} />
                    </div>
                    <span className="logo-text">UpSync</span>
                </Link>
                {!isShare && (
                    <span className="navbar-tagline" aria-hidden="true">
                        Upload · Share · Auto-delete
                    </span>
                )}
            </div>
        </nav>
    );
}
