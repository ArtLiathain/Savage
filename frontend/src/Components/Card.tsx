import React from 'react';

interface CardProps {
  children: React.ReactNode;
}

const Card: React.FC<CardProps> = ({ children }) => {
  return (
    <div className="bg-white rounded-lg shadow-md p-4 overflow-auto min-h-0 h-full w-full">
      {children}
    </div>
  );
};

export default Card;