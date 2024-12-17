import React from 'react';
import { MetricData } from '../../types/DataSnapshot';

interface TableProps {
    data: MetricData[];
}

const MetricDataTable: React.FC<TableProps> = ({ data }) => {
    return (
        <div className="p-4 bg-gray-100">
            <h2 className="text-2xl font-bold text-gray-700 mb-4">Metric Data</h2>
            <div className="overflow-auto">
                <table className="min-w-full border-collapse border border-gray-300 bg-white rounded-md shadow-md">
                    <thead className="bg-gray-200">
                        <tr>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Metric Name</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Device Name</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Device GUID</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Client UTC Time</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Client Timezone (Minutes)</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Server UTC Time</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Server Timezone (Minutes)</th>
                            <th className="border border-gray-300 px-4 py-2 text-left text-sm font-medium text-gray-700">Value</th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.map((metric, index) => (
                            <tr
                                key={index}
                                className={`border border-gray-300 ${
                                    index % 2 === 0 ? 'bg-gray-50' : 'bg-white'
                                }`}
                            >
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700 break-words">{metric.MetricName}</td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700 break-words">{metric.DeviceName}</td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700 break-words">{metric.DeviceGuid}</td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700 break-words">
                                    {new Date(metric.ClientUtcTime).toLocaleString()}
                                </td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700">{metric.ClientTimezoneMinutes}</td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700 break-words">
                                    {new Date(metric.ServerUtcTime).toLocaleString()}
                                </td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700">{metric.ServerTimezoneMinutes}</td>
                                <td className="border border-gray-300 px-4 py-2 text-sm text-gray-700">{metric.Value}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default MetricDataTable;
