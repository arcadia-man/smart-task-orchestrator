import React from 'react';

const CreateScheduler: React.FC = () => {
    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-gray-900">Create Scheduler</h1>
                <p className="mt-2 text-gray-600">
                    Create a new scheduled task
                </p>
            </div>

            <div className="card">
                <p className="text-gray-600">
                    Create scheduler form - to be implemented with:
                </p>
                <ul className="mt-4 space-y-2 text-sm text-gray-600">
                    <li>• Scheduler definition section (name, description, image, command)</li>
                    <li>• Job type selection (immediate, cron, interval)</li>
                    <li>• Cron expression builder</li>
                    <li>• Permission assignment section</li>
                    <li>• Preview of next execution times</li>
                    <li>• Form validation</li>
                </ul>
            </div>
        </div>
    );
};

export default CreateScheduler;