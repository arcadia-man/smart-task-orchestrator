import React from 'react';
import { useParams } from 'react-router-dom';

const SchedulerDetail: React.FC = () => {
    const { id } = useParams<{ id: string }>();

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-gray-900">Scheduler Detail</h1>
                <p className="mt-2 text-gray-600">
                    Viewing scheduler: {id}
                </p>
            </div>

            <div className="card">
                <p className="text-gray-600">
                    Scheduler detail page - to be implemented with:
                </p>
                <ul className="mt-4 space-y-2 text-sm text-gray-600">
                    <li>• Scheduler configuration details</li>
                    <li>• Execution history table</li>
                    <li>• Real-time log viewer</li>
                    <li>• Manual run button</li>
                    <li>• Edit/delete actions</li>
                </ul>
            </div>
        </div>
    );
};

export default SchedulerDetail;