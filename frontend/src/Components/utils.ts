function formatMetricName(metricName: string): string {
    return metricName
      .split('_') // Split the string by underscores
      .map((word, index) => {
        // Capitalize the first letter of each word and make the rest lowercase
        return index === 0
          ? word.charAt(0).toUpperCase() + word.slice(1)
          : word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
      })
      .join(' '); // Join the words back together with spaces
  }

export default formatMetricName