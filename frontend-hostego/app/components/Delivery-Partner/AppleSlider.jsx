import { useState, useRef, useEffect } from "react";

export const AppleStyleSwipeToggle = () => {
  const [isOnline, setIsOnline] = useState(false);
  const sliderRef = useRef(null);
  const knobRef = useRef(null);
  const [position, setPosition] = useState(0);
  const [maxTranslate, setMaxTranslate] = useState(0);

  useEffect(() => {
    if (sliderRef.current && knobRef.current) {
      const sliderWidth = sliderRef.current.offsetWidth;
      const knobWidth = knobRef.current.offsetWidth;
      setMaxTranslate(sliderWidth - knobWidth);
    }
  }, []);

  const handleTouchMove = (e) => {
    if (!sliderRef.current || !knobRef.current) return;

    const slider = sliderRef.current;
    const knob = knobRef.current;
    const touch = e.touches[0];

    let newX = touch.clientX - slider.getBoundingClientRect().left - knob.offsetWidth / 2;
    newX = Math.max(0, Math.min(newX, maxTranslate));

    setPosition(newX);
  };

  const handleTouchEnd = () => {
    if (!sliderRef.current) return;

    const middle = maxTranslate / 2;
    if (position > middle) {
      setIsOnline(true);
      setPosition(maxTranslate); // Move to right
    } else {
      setIsOnline(false);
      setPosition(0); // Move to left
    }
  };

  return (
    <div className="w-full flex justify-center py-6">
      <div
        ref={sliderRef}
        className={`relative w-[90vw] max-w-xl h-16 ${isOnline ? "bg-green-500" : "bg-red-500"
          } rounded-full flex items-center px-2 shadow-xl select-none overflow-hidden`}
      >
        <span
          className={`absolute w-full text-center text-lg font-semibold uppercase transition-all duration-500 ease-out text-gray-200`}
          style={{ transform: `translateX(${isOnline ? maxTranslate - 40 : 16}px)` }}
        >
          {isOnline ? "ONLINE" : "OFFLINE"}
        </span>
        <div
          ref={knobRef}
          className="absolute top-1 left-1 w-14 h-14 bg-white rounded-full shadow-xl transition-transform duration-300 ease-out flex items-center justify-center text-xl font-bold"
          style={{ transform: `translateX(${position}px)` }}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
        >
          {isOnline ? "✔" : "➡"}
        </div>
      </div>
    </div>
  );
};
