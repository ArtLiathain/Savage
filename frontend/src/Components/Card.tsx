import React from 'react';

interface CardProps {
  children: React.ReactNode;
}

const Card: React.FC<CardProps> = ({ children }) => {
  return (
    <div className="bg-white rounded-lg shadow-lg p-6 max-w-4xl ">
      {children}
    </div>
  );
};

export default Card;