# E2EIRC

E2EIRC allows you to create end to end encrpyted chat rooms on Regular (Unmodified) IRC servers and your favorite IRC client.

By default E2EIRC runs a local server on port 6666. After starting up the e2eirc server connect your IRC client to it and it will act as a intermediary between the IRC server and your local client.

# The Gist

Note: This document was written quickly, it assumes you have a good grasp of both AES and RSA encryption.

All messages sent are encrypted by AES256 with a session key unique to each user. The session keys are exchanged using RSA encryption. Each user has a persistent RSA public key. In order for the session key to be exchanged the public key of the requesting use MUST be first approved by the user. Once a key is trusted it will automatically provide the AES key.

# How does it work?

The IRC client connects to the e2eirc process on port 6666. The e2eirc process creates a standard connection to the IRC server of your choice.

The client can then communicate directly with the IRC server through the e2eirc process and join a channel, set a nickname, etc just as it would if it were connected directly.

When a user sends a message the message is intercepted by the e2eirc server and encrypted using the user's own AES key (which changes every session).

The encrypted message is then sent to the IRC server.

On the other end another user running the e2eirc server will recieve the message. It'll attempt to decrypt the message using the last known AES key of the sender. However if no AES key is known, or the AES decryption is unsuccessful it will request the AES key from the sender via private message.

It will also send over it's public key during the handshake. The public key is validated to ensure that it has been seen before. If so the AES encryption key is encrypted with the user's RSA Public Key and sent back to the user.

If the public key has not been seen before it is consiered "untrusted". And the user will be required to approve the key via private message from the user $E2ECtrl.

# Should I use this?

Maybe. Take a look at it, the technology will help you understand how e2e encryption works. It is also a good proof of concept of E2E encryption on both group chats and private chats.

The security, however, is not guarenteed.

Please note that this does not make you anonymous. Your nicknames, ip addresses, etc will all be visable in clear text. What this will do, however, is obscure the content of your mesages. 

# Install

You need to have go installed in order to run this package. It is possible to create pre-built binaries using go, however that isn't the purpose of this project.

Once you have go installed run the following command

```
$ go get github.com/novus0rdo/e2eirc        # download the package
$ go install github.com/novus0rdo/e2eirc    # install it in your gopath bin
```

This will create a binary in your `$GOPATH/bin` directory.

If you have your `$GOPATH/bin` directory in your `PATH` then you can run

```
$ e2eirc -host chat.freenode.net -port 6667  # connects to freenode
```

You will be asked to enter a password. This password will act as your private key encryption key, it is critical that you keep it safe but also remember it.

It should look like the below:

```

███████╗██████╗ ███████╗██╗██████╗  ██████╗
██╔════╝╚════██╗██╔════╝██║██╔══██╗██╔════╝
█████╗   █████╔╝█████╗  ██║██████╔╝██║
██╔══╝  ██╔═══╝ ██╔══╝  ██║██╔══██╗██║
███████╗███████╗███████╗██║██║  ██║╚██████╗
╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═╝ ╚═════╝
Version: 0.0.1-Beta

WARNING: YOU ARE ON A BETA RELEASE!
THE PURPOSE OF THIS RELEASE IS FOR EVALUATION,
SECURITY RESEARCH, AND DEVELOPMENT! THE SECURTIY
OF THIS RELEASE IS NOT GUARENTEED IN ANY WAY.
DO NOT USE THIS VERSION FOR MISSION CRITICAL
COMMUNICATIONS. YOU MAY NOT BE SAFE.

Generating 2048 bit private key. Please wait... Done
Enter a new password for your private key. If you lose it you won't be able to confirm your identity on chat:

```

Enter your password and press enter.

If all looks good, connect your IRC client to 0.0.0.0/6666 (or 0.0.0.0 and port 6666).

If you connected up to freenode join the #e2eirc channel and send a message to try it out!
