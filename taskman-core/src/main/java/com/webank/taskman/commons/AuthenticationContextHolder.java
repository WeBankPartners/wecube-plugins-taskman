package com.webank.taskman.commons;

import java.util.Collection;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

public final class AuthenticationContextHolder {

    private static final ThreadLocal<AuthenticatedUser> currentUser = new InheritableThreadLocal<>();

    public static boolean setAuthenticatedUser(AuthenticatedUser u) {
        if (u == null) {
            return false;
        }

        if (currentUser.get() != null) {
            return false;
        }

        currentUser.set(u);
        return true;
    }

    public static AuthenticatedUser getCurrentUser() {
        return currentUser.get();
    }

    public static void clearCurrentUser() {
        currentUser.remove();
    }

    public static String getCurrentUsername() {
        AuthenticatedUser u = currentUser.get();
        if (u != null) {
            return u.getUsername();
        }

        return null;
    }

    public static Set<String> getCurrentUserRoles() {
        AuthenticatedUser u = currentUser.get();
        if (u != null) {
            return u.getAuthorities();
        }
        
        return Collections.emptySet();
    }

    public static String getCurrentUserRolesToString() {
        AuthenticatedUser u = currentUser.get();
        Set<String> roleSet = new HashSet<>();
        
        //TODO ??
        String roles = "";
        if (u != null && roleSet.size() > 0) {
            roleSet =  u.getAuthorities();
            roles = roleSet.stream().collect(Collectors.joining(","));
        }
        return roles;
    }

    public static class AuthenticatedUser {
        private final String username;
        private String token;
        private Set<String> grantedAuthorities = new HashSet<String>();

        public AuthenticatedUser(String username) {
            super();
            this.username = username;
        }
        
        public AuthenticatedUser(String username, String token, Collection<String> authorities) {
            super();
            this.username = username;
            this.token = token;
            if (authorities != null) {
                for (String a : authorities) {
                    if (a != null) {
                        this.grantedAuthorities.add(a);
                    }
                }
            }
        }

        public String getUsername() {
            return username;
        }

        public Set<String> getAuthorities() {
            return Collections.unmodifiableSet(this.grantedAuthorities);
        }

        public AuthenticatedUser withAuthorities(String... authorities) {
            for (String a : authorities) {
                if (!grantedAuthorities.contains(a)) {
                    grantedAuthorities.add(a);
                }
            }

            return this;
        }

        public String getToken() {
            return token;
        }

        public void setToken(String token) {
            this.token = token;
        }
    }

}