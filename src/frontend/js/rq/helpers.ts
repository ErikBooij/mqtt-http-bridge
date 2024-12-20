import { MaybeAPIError } from '../api/common';

export const unpackMaybeAPIError = <T>([ value, error ]: MaybeAPIError<T>): T => {
    console.log({ value, error });

    if (error !== null) {
        throw new ApiRequestError(error.message);
    }

    return value
}

export class ApiRequestError extends Error {
    errors: string[];

    constructor(message: string | string[]) {
        super(Array.isArray(message) ? message.join(', ') : message);

        this.errors = Array.isArray(message) ? message : [ message ];
    }
}
