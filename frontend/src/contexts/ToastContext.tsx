import React, { createContext, useContext } from 'react';
import { useToast } from '../hooks/useToast';
import ToastContainer from '../components/ToastContainer';

interface ToastContextType {
    success: (title: string, message?: string) => string;
    error: (title: string, message?: string) => string;
    warning: (title: string, message?: string) => string;
    info: (title: string, message?: string) => string;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { toasts, removeToast, success, error, warning, info } = useToast();

    return (
        <ToastContext.Provider value={{ success, error, warning, info }}>
            {children}
            <ToastContainer toasts={toasts} onRemoveToast={removeToast} />
        </ToastContext.Provider>
    );
};

export const useToastContext = () => {
    const context = useContext(ToastContext);
    if (context === undefined) {
        throw new Error('useToastContext must be used within a ToastProvider');
    }
    return context;
};