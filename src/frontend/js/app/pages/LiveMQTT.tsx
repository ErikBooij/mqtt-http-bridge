import React, { useEffect, useState } from 'react'
import { PageTitle } from '../components/PageTitle';
import { useWebSocket } from '../adapters/use-web-socket';
import { ReadyState } from 'react-use-websocket';

type SocketMessage = {
    payload: string;
    topic: string;
    server: string;
    user: string;
    sequence: number;
    timestamp: string;
}

const maxMessages = 1000;

export const LiveMQTT = () => {
    const [ socketMessages, setSocketMessages ] = useState<Record<string, SocketMessage>>({});

    const {
        lastJsonMessage,
        readyState
    } = useWebSocket<SocketMessage>(
        socketURL(),
        {
            shouldReconnect: () => true,
        },
    );

    useEffect(() => {
        if (lastJsonMessage !== null) {
            setSocketMessages((prev) => ( {
                ...prev,
                [ lastJsonMessage.sequence ]: lastJsonMessage
            } ));

            if (Object.keys(socketMessages).length > maxMessages) {
                const newSocketMessages = { ...socketMessages };
                delete newSocketMessages[ Object.keys(socketMessages)[ 0 ] ];
                setSocketMessages(newSocketMessages);
            }
        }
    }, [ lastJsonMessage ]);

    const dateTimeFormatter = new Intl.DateTimeFormat('en-US', {
        minute: 'numeric',
        hour: 'numeric',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        second: 'numeric',
    });

    return (
        <div>
            <PageTitle currentPage="live-mqtt">Live MQTT</PageTitle>
            { readyState !== ReadyState.OPEN && <div>Connection status: { readyState }...</div> }
            { Object.values(socketMessages).length > 0 && (
                <ul role="list" className="divide-y divide-slate-300 flex-col-reverse flex bg-white">
                    { Object.values(socketMessages).map((message) => (
                        <li key={ message.sequence } className="py-5">
                            <div className="flex justify-between gap-x-6 bg-white">
                                <div className="flex min-w-0 gap-x-4">
                                    <div className="min-w-0 flex-auto">
                                        <p className="text-sm/6 font-semibold text-gray-900">{ message.topic }</p>
                                        <p className="mt-1 truncate text-xs/5 text-gray-500"><ServerName
                                            server={ message.server }/></p>
                                    </div>
                                </div>
                                <div className="hidden shrink-0 sm:flex sm:flex-col sm:items-end">
                                    <p className="text-sm/6 text-gray-900">
                                        <time
                                            dateTime={ message.timestamp }
                                        >
                                            { dateTimeFormatter.format(new Date(message.timestamp)) }
                                        </time>
                                    </p>
                                    { message.user &&
                                        <p className="mt-1 text-xs/5 text-gray-500">User: { message.user }</p> }
                                </div>
                            </div>
                            <div>
                                <MessagePayload payload={ message.payload }/>
                            </div>
                        </li>
                    )) }
                </ul>
            ) }
        </div>
    );
}

const socketURL = (): string => {
    const url = new URL('/api/v1/mqtt-socket', window.location.href);

    url.protocol = url.protocol.replace('http', 'ws');

    return url.href;
}

const ServerName = ({ server }: { server: string }) => {
    if (server === '---internal---') {
        return <span
            className="inline-flex items-center gap-x-1.5 rounded-md bg-green-100 px-1.5 py-0.5 text-xs font-medium text-green-700"
        >
            <svg viewBox="0 0 6 6" aria-hidden="true" className="size-1.5 fill-green-500">
                <circle r={ 3 } cx={ 3 } cy={ 3 }/>
            </svg>
            Integrated Broker
      </span>
    }

    return <span
        className="inline-flex items-center gap-x-1.5 rounded-md bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-700"
    >
        <svg viewBox="0 0 6 6" aria-hidden="true" className="size-1.5 fill-blue-500">
            <circle r={ 3 } cx={ 3 } cy={ 3 }/>
        </svg>
        { server }
      </span>
}

const MessagePayload = ({ payload }: { payload: string }) => {
    return <div className="mt-4 text-sm/6 text-gray-900 p-3 bg-white border border-slate-300 rounded-md font-mono">
        { payload }
    </div>
}
