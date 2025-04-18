package com.okblog.post.util;

import java.text.Normalizer;
import java.util.Locale;
import java.util.regex.Pattern;

public class SlugGenerator {
    private static final Pattern NONLATIN = Pattern.compile("[^\\w-]");
    private static final Pattern WHITESPACE = Pattern.compile("[\\s]");
    private static final Pattern MULTIPLE_HYPHENS = Pattern.compile("-+");

    private SlugGenerator() {
        // Private constructor to prevent instantiation
    }

    /**
     * Creates a URL-friendly slug from a string.
     * 
     * @param input The string to convert to a slug
     * @return A URL-friendly slug
     */
    public static String generateSlug(String input) {
        if (input == null || input.isEmpty()) {
            return "";
        }

        String normalized = Normalizer.normalize(input, Normalizer.Form.NFD);
        String noAccents = normalized.replaceAll("\\p{InCombiningDiacriticalMarks}", "");
        String lowercase = noAccents.toLowerCase(Locale.ENGLISH);
        String noWhitespace = WHITESPACE.matcher(lowercase).replaceAll("-");
        String noNonLatin = NONLATIN.matcher(noWhitespace).replaceAll("");
        String noMultipleHyphens = MULTIPLE_HYPHENS.matcher(noNonLatin).replaceAll("-");
        
        // Remove leading/trailing hyphens
        return noMultipleHyphens.replaceAll("^-", "").replaceAll("-$", "");
    }
}