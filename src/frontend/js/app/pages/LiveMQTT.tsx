import React, { useContext, useEffect } from 'react'
import { LayoutContext } from '../components/Layout';

export const LiveMQTT = () => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('live-mqtt');
    })

    return (
        <div>
            <h1>Live MQTT</h1>
            <p>Live MQTT page content...</p>
        </div>
    );
}
