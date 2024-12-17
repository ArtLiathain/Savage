import { Line } from "react-chartjs-2";
import { Chart as ChartJS, CategoryScale, LinearScale, LineElement, PointElement, Title, Tooltip, Legend } from 'chart.js';
import { MetricData } from "../../types/DataSnapshot";

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

interface Dataset {
    label: string;
    data: number[];
    borderColor?: string;
    tension?: number;
    fill?: boolean;
}

interface LineChartProps {
    data: MetricData[] | undefined;
    metricName: string;
    chartTitle: string;
    yAxisLabel: string;
    numberOfDataPoints?: number; // Optional prop to limit number of data points
    datasets?: Dataset[]; // Optional datasets for multi-line charts
}

const LineChart = ({
    data,
    metricName,
    chartTitle,
    yAxisLabel,
    numberOfDataPoints = 40, // Default to the last 40 data points
    datasets,
}: LineChartProps) => {
    if (!data) {
        return <></>;
    }

    const limitedData = data.slice(Math.max(data.length - numberOfDataPoints, 0));

    const labels = limitedData.map((SS) => {
        // Create a Date object from the timestamp string
        const date = new Date(SS.ClientUtcTime);
        // Format it to 24-hour time (HH:mm)
        const formattedTime = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;

        // Get the day of the week (e.g., "Mon", "Tue", etc.)
        const dayOfWeek = date.toLocaleDateString('en-US', { weekday: 'short' });

        // Return a string with the day and time (e.g., "Mon 12:30")
        return `${dayOfWeek} ${formattedTime}`;
    });

    // If no datasets are provided, create one dataset for the given metric
    const chartDatasets = datasets || [{
        label: metricName,
        data: limitedData.map((metric) => metric.Value),
        borderColor: '#4CAF50', // Change color if needed
        tension: 0.4, // Adjust line tension for smoothness
        fill: false, // Make the line graph unfilled
    }];

    return (
        <Line className="w-full"
            data={{
                labels: labels,
                datasets: chartDatasets,
            }}
            options={{
                responsive: true,
                plugins: {
                    title: {
                        display: true,
                        text: chartTitle, // Dynamic chart title
                    },
                    tooltip: {
                        mode: 'index',
                        intersect: false,
                    },
                    legend: {
                        position: 'top',
                    },
                },
                scales: {
                    x: {
                        title: {
                            display: true,
                            text: 'Time (HH:mm)', // X-axis label
                        },
                        type: 'category',
                        offset: true,
                        ticks: {
                            maxRotation: 45,  // Rotate labels for better visibility
                            minRotation: 0,
                        },
                        grid: {
                            display: false,
                        },
                    },
                    y: {
                        title: {
                            display: true,
                            text: yAxisLabel, // Dynamic Y-axis label
                        },
                        suggestedMin: 0,
                        suggestedMax: 100,
                    },
                },
                animation: {
                    duration: 500,  // Animation duration for each data update
                    easing: 'linear', // Smooth linear animation
                    onProgress: (animation) => {
                        const chartInstance = animation.chart;
                        const ctx = chartInstance.ctx;
                        const xScale = chartInstance.options.scales?.x;

                        if (xScale) {
                            ctx.save();
                            // Update the x-axis range as the chart scrolls
                            xScale.min = chartInstance.chartArea.left;
                            xScale.max = chartInstance.chartArea.right;
                            ctx.restore();
                        }
                    }
                },
                layout: {
                    padding: {
                        top: 10,
                        left: 10,
                        right: 10,
                        bottom: 10,
                    },
                },
            }}
        />
    );
};

export default LineChart;
