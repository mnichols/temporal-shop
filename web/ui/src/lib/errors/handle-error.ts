import { browser } from '$app/env';
import { networkError } from '$lib/stores/error';
import type { APIErrorResponse } from './request-from-api';
import { routeForLoginPage } from './route-for';

// This will eventually be expanded on.
export const handleError = (
    error: any,
    errors = networkError,
    isBrowser = browser,
): void => {
    if (isUnauthorized(error) && isBrowser) {
        window.location.assign(routeForLoginPage(error?.message));
        return;
    }

    if (isForbidden(error) && isBrowser) {
        window.location.assign(routeForLoginPage(error?.message));
        return;
    }
};

export const handleUnauthorizedOrForbiddenError = (
    error: APIErrorResponse,
    isBrowser = browser,
): void => {
    const msg = `${error?.status} ${error?.body?.message}`;

    if (isUnauthorized(error) && isBrowser) {
        window.location.assign(routeForLoginPage(msg));
        return;
    }

    if (isForbidden(error) && isBrowser) {
        window.location.assign(routeForLoginPage(msg));
        return;
    }
};

export const isUnauthorized = (error: any): boolean => {
    return error?.statusCode === 401 || error?.status === 401;
};

export const isForbidden = (error: any): boolean => {
    return error?.statusCode === 403 || error?.status === 403;
};