import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ToastProvider } from '@/context/ToastContext';
import { ErrorBoundary } from '@/components/ErrorBoundary';
import ToastContainer from '@/components/ToastContainer';
import Navbar from '@/components/Navbar';
import UploadPage from '@/pages/UploadPage';
import SharePage from '@/pages/SharePage';
import './index.css';

function Footer() {
    return (
        <footer className="footer" aria-label="Site footer">
            Built with ❤️ · Files auto-delete after expiry · <span>UpSync</span>
        </footer>
    );
}

export default function App() {
    return (
        <ErrorBoundary>
            <ToastProvider>
                <BrowserRouter>
                    <div className="app-wrapper">
                        <div className="bg-orbs" aria-hidden="true" />

                        <Routes>
                            {/* Share/download page — full-screen layout, no navbar */}
                            <Route
                                path="/share/:id"
                                element={
                                    <ErrorBoundary>
                                        <SharePage />
                                    </ErrorBoundary>
                                }
                            />

                            {/* Upload page — standard layout with navbar + footer */}
                            <Route
                                path="/*"
                                element={
                                    <>
                                        <Navbar />
                                        <ErrorBoundary>
                                            <UploadPage />
                                        </ErrorBoundary>
                                        <Footer />
                                    </>
                                }
                            />
                        </Routes>

                        {/* Global toast notifications */}
                        <ToastContainer />
                    </div>
                </BrowserRouter>
            </ToastProvider>
        </ErrorBoundary>
    );
}
