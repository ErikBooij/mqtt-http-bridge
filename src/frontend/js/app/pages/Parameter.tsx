import React, { useContext, useEffect } from 'react'
import { LayoutContext } from '../components/Layout';

export const Parameter = () => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('parameters');
    })

    return (
        <div>
            <h1>Parameter</h1>
            <p>Parameter page content...</p>
        </div>
    );
}
