import React from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Subscriptions } from './pages/Subscriptions';
import { Layout } from './components/Layout';
import { Subscription } from './pages/Subscription';
import { Parameters } from './pages/Parameters';
import { Parameter } from './pages/Parameter';
import { LiveMQTT } from './pages/LiveMQTT';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

export const App = () => {
    return (
        <React.StrictMode>
            <QueryClientProvider client={ queryClient }>
                <BrowserRouter>
                    <Routes>
                        <Route path="/" element={ <Layout/> }>
                            {/* Overview page (subscription list) */ }
                            <Route path="" element={ <Subscriptions/> }/>
                            {/* Individual subscription page (as edit form) */ }
                            <Route path="subscriptions/:id" element={ <Subscription mode="edit"/> }/>
                            {/* New subscription page (same as edit form, with no prefill, and different save action */ }
                            <Route path="new-subscription" element={ <Subscription mode="new"/> }/>
                            {/* Copy subscription page (same as new form, with prefill from existing subscription) */ }
                            <Route path="copy-subscription/:id" element={ <Subscription mode="copy"/> }/>
                            {/* Overview page (parameter list) */ }
                            <Route path="parameters" element={ <Parameters/> }/>
                            {/* Individual parameter page (as edit form) */ }
                            <Route path="parameter/:id" element={ <Parameter/> }/>
                            {/* New parameter page (same as edit form, with no prefill, and different save action) */ }
                            <Route path="new-parameter" element={ <Parameter/> }/>
                            {/* Page with live MQTT data from all connected servers */ }
                            <Route path="live-mqtt" element={ <LiveMQTT/> }/>
                        </Route>
                    </Routes>
                </BrowserRouter>
            </QueryClientProvider>
        </React.StrictMode>
    );
}
