'use strict';

const gulp         = require("gulp");
const sass         = require("gulp-sass");
const sourcemaps   = require('gulp-sourcemaps');
const autoprefixer = require('gulp-autoprefixer');
const babel        = require('gulp-babel');

const dir = {
    src: 'resources',
    dist: 'public'
};

const options = {
    sass : { outputStyle: 'compressed' },
    es6 : { presets : ['es2015'] },
    autoprefixer : {
        browsers : ['last 5 versions'],
        cascade : false
    }
};

gulp.task('sass', () => {
    gulp.src(`${dir.src}/scss/**/*.scss`)
        .pipe(sourcemaps.init())
        .pipe(sass(options.sass).on('error', sass.logError))
        .pipe(autoprefixer(options.autoprefixer))
        .pipe(sourcemaps.write('./'))
        .pipe(gulp.dest(`${dir.dist}/css`));
});

gulp.task('es6', () => {
    gulp.src(`${dir.src}/es6/*.js`)
        .pipe(babel(options.es6))
        .pipe(gulp.dest(`${dir.dist}/js`));
});

gulp.task('watch', () => {
    gulp.watch(`${dir.src}/scss/**/*.scss`, ['sass']);
    gulp.watch(`${dir.src}/es6/*.js`, ['es6']);
})

gulp.task('default', ['sass', 'es6', 'watch']);
