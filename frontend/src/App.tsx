import Linechart from './Components/Linechart'
import { DataSnapshot } from '../types/DataSnapshot'
import Card from './Components/Card'
import { useEffect, useState } from 'react'

function App() {
  const [metrics, setMetrics] =  useState<DataSnapshot[] | undefined>(undefined)
  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const response = await fetch("http://localhost:8080");
        if (!response.ok) {
          throw new Error('Failed to fetch data');
        }
        const data = await response.json(); // This resolves to the actual JSON data
        setMetrics(data); // Now you pass the resolved data (array of DataSnapshot)
      } catch (error) {
        console.error('Error fetching metrics:', error);
      }
    };

    fetchMetrics()
  }, [])

  return (
    <div className='bg-blue-500 min-h-screen'>
      <Card>
        <Linechart data={metrics} ></Linechart>
      </Card>
    </div>
  )
}

export default App
