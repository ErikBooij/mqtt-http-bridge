import vendorUseWebSocket from 'react-use-websocket';

// @ts-expect-error - This is a hack to get around the fact that the type definition for react-use-websocket is incorrect.
export const useWebSocket = vendorUseWebSocket.default! as typeof vendorUseWebSocket;
