import React from 'react';
import { GaugeComponent } from 'react-gauge-component';
import { MetricData } from '../../types/DataSnapshot';
import formatMetricName from './utils';



// Define the GaugeComponent that accepts MetricData
const GaugeDisplay: React.FC<{ metric: MetricData }> = ({ metric }) => {
    // Define the minimum and maximum values for the gauge (you can adjust these as needed)

    return (
        <div className='flex flex-col justify-center items-center w-fullh-full'>
            <p>{formatMetricName(metric.MetricName)}</p>

            <GaugeComponent
                type="semicircle"
                arc={{
                    colorArray: ['#00FF15', '#FF2121'],
                    padding: 0.02,
                    subArcs:
                        [
                            { limit: 40 },
                            { limit: 60 },
                            { limit: 70 },
                            {},
                            {},
                            {},
                            {}
                        ]
                }}
                pointer={{ type: "blob", animationDelay: 0 }}
                value={metric.Value}
            />
        </div>
    );
};
export default GaugeDisplay