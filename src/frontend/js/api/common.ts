import { ZodError } from 'zod';

export type APIError = {
    statusCode: number;
    message: string | string[];
}

export type MaybeAPIError<T> = [ T, null ] | [ null, APIError ];
export type AsyncMaybeAPIError<T> = Promise<MaybeAPIError<T>>;

export const apiErrorFromResponse = async (response: Response): Promise<APIError> => {
    try {
        const body = await response.json();

        if (isApiError(body)) {
            return { statusCode: response.status, message: body.error };
        }
    } catch (e) {
    }

    return { statusCode: response.status, message: 'Unknown error' };
}

type errorBodyType = {
    error: string | string[];
}

const isApiError = (body: unknown): body is errorBodyType => {
    return body !== null && typeof body === 'object' && 'error' in body && (typeof body.error === 'string' || Array.isArray(body.error));
}

export const responseParseError = (error: ZodError): APIError => {
    return { statusCode: 0, message: error.toString() };
}
