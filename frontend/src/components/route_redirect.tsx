import { Navigate } from 'react-router-dom';

export function RouteRedirect() {
    return <Navigate to="/admin" replace />;
}