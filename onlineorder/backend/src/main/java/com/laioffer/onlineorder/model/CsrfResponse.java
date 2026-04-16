package com.laioffer.onlineorder.model;

public record CsrfResponse(String token, String headerName, String parameterName) {
}
