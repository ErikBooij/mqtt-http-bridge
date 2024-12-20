import React, { useContext, useEffect } from 'react'
import { LayoutContext } from '../components/Layout';

export const Parameters = () => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('parameters');
    })

    return (
        <div>
            <h1>Parameters</h1>
            <p>Parameters page content...</p>
        </div>
    );
}
