import { DataSnapshot } from '../../types/DataSnapshot'
import { Line } from "react-chartjs-2"

import { Chart as ChartJS, CategoryScale, LinearScale, LineElement, PointElement, Title, Tooltip, Legend } from 'chart.js';

// Register necessary components of Chart.js
ChartJS.register(
    CategoryScale,   // For x-axis "category" scale
    LinearScale,     // For y-axis scale
    LineElement,     // For drawing the line
    PointElement,    // For drawing the points on the line
    Title,           // For title
    Tooltip,         // For tooltips
    Legend           // For legend
);

interface ExampleComponentProps {
    data: DataSnapshot[] | undefined;
}

const Linechart = ({ data }: ExampleComponentProps) => {
    if (!data) {
        return (<></>)
    }
    return (
        <Line data={{
            labels: data.map((SS) => {
                // Create a Date object from the timestamp string
                const date = new Date(SS.Timestamp);

                // Format it to 24-hour time (HH:mm)
                return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
            }),
            datasets: [{
                label: 'TOTAL RAM USAGE',
                data: data.flatMap((snapshot) => {
                    const totalMetric = snapshot.Metrics.find((temp) => temp.Name === "Total");
                    return totalMetric ? totalMetric.Value : null;  // Handle missing "Total" metric
                }).filter((value) => value !== null),

            }]

        }}></Line>
    )
}

export default Linechart