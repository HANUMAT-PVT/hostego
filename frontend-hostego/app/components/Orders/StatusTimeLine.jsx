import React from 'react'

const StatusTimeLine = ({ ORDER_STATUSES, activeOrder }) => {
  
    const getStatusStep = (status) => {
        return ORDER_STATUSES.findIndex(s => s.id === status);
    }
    return (
        <div div className="bg-white rounded-xl shadow-sm p-4 mb-20" >
            <h4 className="font-semibold mb-4">Order Progress</h4>
            <div className="relative">
      {/* Background Vertical Line */}
      <div
        className="absolute left-[15px] top-[24px] w-[2px] bg-gray-100"
        style={{
          height: `${(ORDER_STATUSES.length - 1) * 56}px`  // Adjusted height calculation
        }}
      />

      {/* Colored Progress Line */}
      <div
        className="absolute left-[15px] top-[24px] w-[2px] transition-all duration-300"
        style={{
          height: `${getStatusStep(activeOrder.order_status) * 56}px`,
          backgroundColor: 'var(--primary-color)'
        }}
      />

      {/* Status Points */}
      <div className="relative">
        {ORDER_STATUSES.map((status, index) => {
          const isCompleted = getStatusStep(activeOrder.order_status) >= index;
          const isCurrent = getStatusStep(activeOrder.order_status) === index;
          const StatusIcon = status.icon;

          return (
            <div
              key={status.id}
              className="flex items-center gap-4 relative h-14" // Fixed height for each status

              style={{ cursor: 'pointer' }}
            >
              {/* Status Circle */}
              <div
                style={{
                  backgroundColor: isCompleted ? status.color : '#fff',
                  borderColor: isCompleted ? status.color : '#e5e7eb',
                }}
                className={`
                  w-8 h-8 rounded-full flex items-center justify-center 
                  border-2 relative z-10 bg-white transition-all duration-300
                `}
              >
                <StatusIcon
                  className={`w-4 h-4 ${isCompleted ? 'text-white' : 'text-gray-400'
                    }`}
                />
              </div>

              {/* Status Text */}
              <div className="flex flex-col">
                <span className={`font-medium ${isCompleted ? 'text-gray-900' : 'text-gray-400'
                  }`}>
                  {status.label}
                </span>
                {isCurrent && (
                  <span className="text-sm text-[var(--primary-color)]">
                    Current
                  </span>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  </div>
  )
}

export default StatusTimeLine
