package com.webank.taskman.authclient.filter;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jws;

public interface JwtSsoTokenParser {
    Jws<Claims> parseJwt(String token);
}
