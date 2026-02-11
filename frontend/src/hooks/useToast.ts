import { useState, useCallback } from 'react';
import { ToastData, ToastType } from '../components/ToastContainer';

export const useToast = () => {
    const [toasts, setToasts] = useState<ToastData[]>([]);

    const addToast = useCallback((
        type: ToastType,
        title: string,
        message?: string,
        duration?: number
    ) => {
        const id = Math.random().toString(36).substr(2, 9);
        const newToast: ToastData = {
            id,
            type,
            title,
            message,
            duration,
        };

        setToasts((prev) => [...prev, newToast]);
        return id;
    }, []);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
    }, []);

    const success = useCallback((title: string, message?: string) => {
        return addToast('success', title, message);
    }, [addToast]);

    const error = useCallback((title: string, message?: string) => {
        return addToast('error', title, message, 7000); // Longer duration for errors
    }, [addToast]);

    const warning = useCallback((title: string, message?: string) => {
        return addToast('warning', title, message);
    }, [addToast]);

    const info = useCallback((title: string, message?: string) => {
        return addToast('info', title, message);
    }, [addToast]);

    return {
        toasts,
        addToast,
        removeToast,
        success,
        error,
        warning,
        info,
    };
};