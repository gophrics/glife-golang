import Register


def main():
    success = Register.register()
    if success:
        Register.get_me()
        Register.get_user()

if __name__ == "__main__":
    main()
